package main

import (
	"encoding/hex"
	"io/ioutil"
	"net"
)

func ProcessRequest(addr *net.UDPAddr, buf []byte, s *Server) {
	resolveNormal(s, addr, buf)
	//resolvedQuery := DNSOverHTTPSRequest(base64.StdEncoding.EncodeToString(buf), s.config) //request DNS record over HTTPS
	////TODO: Fallback to normal DNS query? e.g. in case of Captive Portal
	//if resolvedQuery != nil {
	//	_, err := s.conn.WriteToUDP(resolvedQuery, addr)
	//	CheckError(err)
	//}
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

func resolveNormal(s *Server, addr *net.UDPAddr, buf []byte) {
	b := make([]byte, 1024)
	udpAddr, err := net.ResolveUDPAddr("udp", "1.1.1.1:53")
	CheckError(err)
	outgoing, err := net.DialUDP("udp", nil, udpAddr)
	CheckError(err)
	go func(addr *net.UDPAddr) {
		for {
			n1, _, e := outgoing.ReadFromUDP(b)
			CheckError(e)
			_, e = s.conn.WriteToUDP(b[0:n1], addr)
			CheckError(e)
		}
	}(addr)
	_, e := outgoing.Write(buf)
	CheckError(e)
}
