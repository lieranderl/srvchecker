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

			for _, eku := range cert.ExtKeyUsage {
				if eku == x509.ExtKeyUsageServerAuth {
					c.ExtKeyUsage = append(c.ExtKeyUsage, "TLS Web Server Authentication")
				}
				if eku == x509.ExtKeyUsageClientAuth {
					c.ExtKeyUsage = append(c.ExtKeyUsage, "TLS Web Client Authentication")
				}
			}

			c.Txt, _ = certinfo.CertificateText(cert)
			ps.CertsChain = append(ps.CertsChain, c)

		}

	}
}
