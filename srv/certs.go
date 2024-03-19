package srv

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/grantae/certinfo"
)

type Cert struct {
	Txt         string
	Cn          string
	Subject     string
	San         string
	KeyUsage    []string
	ExtKeyUsage []string
	Issuer      string
	NotBefore   string
	NotAfter    string
}
type CertsChain []*Cert

func getCert(ip, port string) ([]*x509.Certificate, error) {
	conf := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 2 * time.Second}, "tcp", net.JoinHostPort(ip, port), conf)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.ConnectionState().PeerCertificates, nil
}

func keyUsageToStrings(keyUsage x509.KeyUsage) []string {
	usages := []string{}

	if keyUsage&x509.KeyUsageDigitalSignature != 0 {
		usages = append(usages, "Digital Signature")
	}
	if keyUsage&x509.KeyUsageContentCommitment != 0 {
		usages = append(usages, "Content Commitment")
	}
	if keyUsage&x509.KeyUsageKeyEncipherment != 0 {
		usages = append(usages, "Key Encipherment")
	}
	if keyUsage&x509.KeyUsageDataEncipherment != 0 {
		usages = append(usages, "Data Encipherment")
	}
	if keyUsage&x509.KeyUsageKeyAgreement != 0 {
		usages = append(usages, "Key Agreement")
	}
	if keyUsage&x509.KeyUsageCertSign != 0 {
		usages = append(usages, "Certificate Sign")
	}
	if keyUsage&x509.KeyUsageCRLSign != 0 {
		usages = append(usages, "CRL Sign")
	}
	if keyUsage&x509.KeyUsageEncipherOnly != 0 {
		usages = append(usages, "Encipher Only")
	}
	if keyUsage&x509.KeyUsageDecipherOnly != 0 {
		usages = append(usages, "Decipher Only")
	}

	return usages
}

func extKeyUsageToStrings(extKeyUsages []x509.ExtKeyUsage) []string {
	var usageStrings []string

	for _, usage := range extKeyUsages {
		switch usage {
		case x509.ExtKeyUsageAny:
			usageStrings = append(usageStrings, "Any")
		case x509.ExtKeyUsageServerAuth:
			usageStrings = append(usageStrings, "Server Authentication")
		case x509.ExtKeyUsageClientAuth:
			usageStrings = append(usageStrings, "Client Authentication")
		case x509.ExtKeyUsageCodeSigning:
			usageStrings = append(usageStrings, "Code Signing")
		case x509.ExtKeyUsageEmailProtection:
			usageStrings = append(usageStrings, "Email Protection")
		case x509.ExtKeyUsageIPSECEndSystem:
			usageStrings = append(usageStrings, "IPSEC End System")
		case x509.ExtKeyUsageIPSECTunnel:
			usageStrings = append(usageStrings, "IPSEC Tunnel")
		case x509.ExtKeyUsageIPSECUser:
			usageStrings = append(usageStrings, "IPSEC User")
		case x509.ExtKeyUsageTimeStamping:
			usageStrings = append(usageStrings, "Time Stamping")
		case x509.ExtKeyUsageOCSPSigning:
			usageStrings = append(usageStrings, "OCSP Signing")
		case x509.ExtKeyUsageMicrosoftServerGatedCrypto:
			usageStrings = append(usageStrings, "Microsoft Server Gated Crypto")
		case x509.ExtKeyUsageNetscapeServerGatedCrypto:
			usageStrings = append(usageStrings, "Netscape Server Gated Crypto")
		case x509.ExtKeyUsageMicrosoftCommercialCodeSigning:
			usageStrings = append(usageStrings, "Microsoft Commercial Code Signing")
		case x509.ExtKeyUsageMicrosoftKernelCodeSigning:
			usageStrings = append(usageStrings, "Microsoft Kernel Code Signing")
		default:
			usageStrings = append(usageStrings, fmt.Sprintf("Unknown(%d)", usage))
		}
	}

	return usageStrings
}

func (ps *DiscoveredSrvRow) Connect_cert(ip string, port string) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		ps.IsOpened = false
	}
	if err == nil {
		defer conn.Close()
		ps.IsOpened = true
		certs, err := getCert(ip, fmt.Sprint(port))
		// check error
		if err != nil {
			ps.CertsChain = nil
		}

		for _, cert := range certs {
			c := new(Cert)
			c.Cn = cert.Subject.CommonName
			c.Issuer = cert.Issuer.CommonName
			c.Subject = cert.Subject.String()
			c.San = strings.Join(cert.DNSNames, ", ")
			c.NotBefore = cert.NotBefore.String()
			c.NotAfter = cert.NotAfter.String()
			c.KeyUsage = keyUsageToStrings(cert.KeyUsage)
			c.ExtKeyUsage = extKeyUsageToStrings(cert.ExtKeyUsage)
			c.Txt, _ = certinfo.CertificateText(cert)
			ps.CertsChain = append(ps.CertsChain, c)

		}

	}
}
