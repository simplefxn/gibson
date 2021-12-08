package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"sync"
)

var onceByKeys = sync.Map{}
var combinedTLSCertificates = sync.Map{}
var tlsCertificates = sync.Map{}
var certPools = sync.Map{}

// Fixed upstream in https://github.com/golang/go/issues/13385
func newTLSConfig(minVersion uint16) *tls.Config {

	ciphers := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	}

	if minVersion < tls.VersionTLS12 {
		ciphers = append(ciphers, []uint16{
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		}...)
	}

	return &tls.Config{
		MinVersion:   minVersion,
		CipherSuites: ciphers,
	}
}

func tlsCertificatesIdentifier(tokens ...string) string {
	return strings.Join(tokens, ";")
}

func doLoadAndCombineTLSCertificates(ca, cert, key string) error {
	combinedTLSIdentifier := tlsCertificatesIdentifier(ca, cert, key)

	// Read CA certificates chain
	caB, err := os.ReadFile(ca)
	if err != nil {
		return fmt.Errorf("failed to read ca file: %s", ca)
	}

	// Read server certificate
	certB, err := os.ReadFile(cert)
	if err != nil {
		return fmt.Errorf("failed to read server cert file: %s", cert)
	}

	// Read server key file
	keyB, err := os.ReadFile(key)
	if err != nil {
		return fmt.Errorf("failed to read key file: %s", key)
	}

	// Load CA, server cert and key.
	var certificate []tls.Certificate
	crt, err := tls.X509KeyPair(append(certB, caB...), keyB)
	if err != nil {
		return fmt.Errorf("failed to load and merge tls certificate with CA, ca %s, cert %s, key: %s", ca, cert, key)
	}

	certificate = []tls.Certificate{crt}

	combinedTLSCertificates.Store(combinedTLSIdentifier, &certificate)

	return nil
}

func combineAndLoadTLSCertificates(ca, cert, key string) (*[]tls.Certificate, error) {
	combinedTLSIdentifier := tlsCertificatesIdentifier(ca, cert, key)
	once, _ := onceByKeys.LoadOrStore(combinedTLSIdentifier, &sync.Once{})

	var err error
	once.(*sync.Once).Do(func() {
		err = doLoadAndCombineTLSCertificates(ca, cert, key)
	})

	if err != nil {
		return nil, err
	}

	result, ok := combinedTLSCertificates.Load(combinedTLSIdentifier)

	if !ok {
		return nil, fmt.Errorf("Cannot find loaded tls certificate chain with ca: %s, cert: %s, key: %s", ca, cert, key)
	}

	return result.(*[]tls.Certificate), nil
}

func doLoadTLSCertificate(cert, key string) error {
	tlsIdentifier := tlsCertificatesIdentifier(cert, key)

	var certificate []tls.Certificate
	// Load the server cert and key.
	crt, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return fmt.Errorf("failed to load tls certificate, cert %s, key: %s", cert, key)
	}

	certificate = []tls.Certificate{crt}

	tlsCertificates.Store(tlsIdentifier, &certificate)

	return nil
}

func loadTLSCertificate(cert, key string) (*[]tls.Certificate, error) {
	tlsIdentifier := tlsCertificatesIdentifier(cert, key)
	once, _ := onceByKeys.LoadOrStore(tlsIdentifier, &sync.Once{})

	var err error
	once.(*sync.Once).Do(func() {
		err = doLoadTLSCertificate(cert, key)
	})

	if err != nil {
		return nil, err
	}

	result, ok := tlsCertificates.Load(tlsIdentifier)

	if !ok {
		return nil, fmt.Errorf("Cannot find loaded tls certificate with cert: %s, key%s", cert, key)
	}

	return result.(*[]tls.Certificate), nil
}

func doLoadx509CertPool(ca string) error {
	b, err := os.ReadFile(ca)
	if err != nil {
		return fmt.Errorf("failed to read ca file: %s", ca)
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return fmt.Errorf("failed to append certificates")
	}

	certPools.Store(ca, cp)

	return nil
}

func loadx509CertPool(ca string) (*x509.CertPool, error) {
	once, _ := onceByKeys.LoadOrStore(ca, &sync.Once{})

	var err error
	once.(*sync.Once).Do(func() {
		err = doLoadx509CertPool(ca)
	})
	if err != nil {
		return nil, err
	}

	result, ok := certPools.Load(ca)

	if !ok {
		return nil, fmt.Errorf("Cannot find loaded x509 cert pool for ca: %s", ca)
	}

	return result.(*x509.CertPool), nil
}
