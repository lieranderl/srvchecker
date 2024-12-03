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

// Cert represents a parsed X.509 certificate
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

// CertsChain is a collection of Cert objects
type CertsChain []*Cert

// getCert fetches the X.509 certificates from a given IP and port
func getCert(ip, port string) ([]*x509.Certificate, error) {
	conf := &tls.Config{InsecureSkipVerify: true}
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 2 * time.Second}, "tcp", net.JoinHostPort(ip, port), conf)
	if err != nil {
		return nil, fmt.Errorf("failed to establish TLS connection: %v", err)
	}
	defer conn.Close()
	return conn.ConnectionState().PeerCertificates, nil
}

// parseCertificates converts raw X.509 certificates to Cert objects
func parseCertificates(certs []*x509.Certificate) CertsChain {
	var certsChain CertsChain
	for _, cert := range certs {
		certsChain = append(certsChain, &Cert{
			Cn:          cert.Subject.CommonName,
			Subject:     cert.Subject.String(),
			San:         strings.Join(cert.DNSNames, ", "),
			KeyUsage:    keyUsageToStrings(cert.KeyUsage),
			ExtKeyUsage: extKeyUsageToStrings(cert.ExtKeyUsage),
			Issuer:      cert.Issuer.CommonName,
			NotBefore:   cert.NotBefore.String(),
			NotAfter:    cert.NotAfter.String(),
			Txt:         parseCertInfo(cert),
		})
	}
	return certsChain
}

// parseCertInfo generates a detailed certificate info text
func parseCertInfo(cert *x509.Certificate) string {
	txt, err := certinfo.CertificateText(cert)
	if err != nil {
		return "Failed to parse certificate details"
	}
	return txt
}

// keyUsageToStrings converts key usages to string representations
func keyUsageToStrings(keyUsage x509.KeyUsage) []string {
	usages := []string{}

	usageMap := map[x509.KeyUsage]string{
		x509.KeyUsageDigitalSignature:  "Digital Signature",
		x509.KeyUsageContentCommitment: "Content Commitment",
		x509.KeyUsageKeyEncipherment:   "Key Encipherment",
		x509.KeyUsageDataEncipherment:  "Data Encipherment",
		x509.KeyUsageKeyAgreement:      "Key Agreement",
		x509.KeyUsageCertSign:          "Certificate Sign",
		x509.KeyUsageCRLSign:           "CRL Sign",
		x509.KeyUsageEncipherOnly:      "Encipher Only",
		x509.KeyUsageDecipherOnly:      "Decipher Only",
	}

	for k, v := range usageMap {
		if keyUsage&k != 0 {
			usages = append(usages, v)
		}
	}

	return usages
}

// extKeyUsageToStrings converts extended key usages to string representations
func extKeyUsageToStrings(extKeyUsages []x509.ExtKeyUsage) []string {
	var usageStrings []string

	usageMap := map[x509.ExtKeyUsage]string{
		x509.ExtKeyUsageAny:                            "Any",
		x509.ExtKeyUsageServerAuth:                     "Server Authentication",
		x509.ExtKeyUsageClientAuth:                     "Client Authentication",
		x509.ExtKeyUsageCodeSigning:                    "Code Signing",
		x509.ExtKeyUsageEmailProtection:                "Email Protection",
		x509.ExtKeyUsageIPSECEndSystem:                 "IPSEC End System",
		x509.ExtKeyUsageIPSECTunnel:                    "IPSEC Tunnel",
		x509.ExtKeyUsageIPSECUser:                      "IPSEC User",
		x509.ExtKeyUsageTimeStamping:                   "Time Stamping",
		x509.ExtKeyUsageOCSPSigning:                    "OCSP Signing",
		x509.ExtKeyUsageMicrosoftServerGatedCrypto:     "Microsoft Server Gated Crypto",
		x509.ExtKeyUsageNetscapeServerGatedCrypto:      "Netscape Server Gated Crypto",
		x509.ExtKeyUsageMicrosoftCommercialCodeSigning: "Microsoft Commercial Code Signing",
		x509.ExtKeyUsageMicrosoftKernelCodeSigning:     "Microsoft Kernel Code Signing",
	}

	for _, usage := range extKeyUsages {
		if usageName, found := usageMap[usage]; found {
			usageStrings = append(usageStrings, usageName)
		} else {
			usageStrings = append(usageStrings, fmt.Sprintf("Unknown(%d)", usage))
		}
	}

	return usageStrings
}

// Connect_cert checks connectivity and retrieves certificates for a discovered SRV record
func (ps *DiscoveredSrvRow) Connect_cert(ip, port string) {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, port), timeout)
	if err != nil {
		ps.IsOpened = false
		return
	}
	defer conn.Close()

	ps.IsOpened = true

	// Fetch and parse certificates
	certs, err := getCert(ip, port)
	if err != nil {
		ps.CertsChain = nil
		fmt.Printf("Error fetching certificates for %s:%s: %v\n", ip, port, err)
		return
	}

	ps.CertsChain = parseCertificates(certs)
}
