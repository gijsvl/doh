package main

import (
	"bytes"
	"crypto"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"net"
	"net/http"
)

type HttpPKPClient interface {
	Get(url string) (*http.Response, error)
	Post(url string, body io.Reader) (*http.Response, error)
}

type Client struct {
	fingerprint []byte
}

type Dialer func(network, addr string) (net.Conn, error)

func NewClient(fingerprint []byte) HttpPKPClient {
	return Client{fingerprint}
}

func computeFingerprint(pubKey crypto.PublicKey) ([32]byte, error) {
	der, err := x509.MarshalPKIXPublicKey(pubKey)
	return sha256.Sum256(der), err
}

func (c Client) makeDialer(fingerprint []byte, skipCAVerification bool) Dialer {
	return func(network, addr string) (net.Conn, error) {
		conn, err := tls.Dial(network, addr, &tls.Config{InsecureSkipVerify: skipCAVerification})
		if err != nil {
			return conn, err
		}
		connState := conn.ConnectionState()
		keyPinValid := false
		for _, peerCert := range connState.PeerCertificates {
			hash, err := computeFingerprint(peerCert.PublicKey)
			if err != nil {
				log.Fatal(err)
			}
			if bytes.Compare(hash[0:], fingerprint) == 0 {
				keyPinValid = true
			}
		}
		if keyPinValid == false {
			println("invalid pin")
		}
		return conn, nil
	}
}

func (c Client) Get(url string) (*http.Response, error) {
	client := &http.Client{}
	client.Transport = &http.Transport{
		DialTLS: c.makeDialer(c.fingerprint, false),
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("accept", "application/dns-message")
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

func (c Client) Post(url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	client.Transport = &http.Transport{
		DialTLS: c.makeDialer(c.fingerprint, false),
	}
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("content-type", "application/dns-message")
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}
