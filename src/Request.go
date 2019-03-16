package main

import (
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net"
)

func ProcessRequest(addr *net.UDPAddr, buf []byte, s *Server) {
	resolvedQuery := DNSOverHTTPSRequest(base64.StdEncoding.EncodeToString(buf), s.config) //request DNS record over HTTPS
	//TODO: Fallback to normal DNS query? e.g. in case of Captive Portal
	if resolvedQuery != nil {
		_, err := s.conn.WriteToUDP(resolvedQuery, addr)
		CheckError(err)
	}
}

func DNSOverHTTPSRequest(record string, config *Config) []byte {
	//QUERY OVER HTTPS
	bytes, _ := hex.DecodeString(config.ChosenEndpoint.Fingerprint)
	client := NewClient(bytes)
	res, err := client.Get(config.ChosenEndpoint.Url + record)
	if err != nil {
		//TODO: fallback to other endpoint?
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	CheckError(err)

	return body
}
