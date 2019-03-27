package main

import (
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"net"
)

func ProcessRequest(addr *net.UDPAddr, buf []byte, s *Server) {
	resolvedQuery := DNSOverHTTPSRequest(base64.StdEncoding.EncodeToString(buf), s.config) //request DNS record over HTTPS
	if resolvedQuery != nil {
		_, err := s.conn.WriteToUDP(resolvedQuery, addr)
		if !CheckError(err) {
			return
		}
	} else {
		Notification("Fallback to normal DNS, because something went wrong with DoH. You lost privacy!")
		resolveNormal(s, addr, buf)
	}
}

func DNSOverHTTPSRequest(record string, config *Config) []byte {
	//QUERY OVER HTTPS
	bytes, _ := hex.DecodeString(config.ChosenEndpoint.Fingerprint)
	client := NewClient(bytes)
	res, err := client.Get(config.ChosenEndpoint.Url + record)
	if err != nil {
		for _, endpoint := range config.Endpoints {
			res, err = client.Get(endpoint.Url + record)
			if err != nil {
				continue
			} else {
				Notification("Fallback to other endpoint: " + endpoint.Name)
			}
		}
		if err != nil {
			return nil
		}
	}
	body, err := ioutil.ReadAll(res.Body)
	if !CheckError(err) {
		return nil
	}

	return body
}

func resolveNormal(s *Server, addr *net.UDPAddr, buf []byte) {
	b := make([]byte, 1024)
	udpAddr, err := net.ResolveUDPAddr("udp", "1.1.1.1:53")
	if err != nil {
		return
	}
	outgoing, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return
	}
	go func(addr *net.UDPAddr) {
		for {
			n1, _, e := outgoing.ReadFromUDP(b)
			if !CheckError(e) {
				return
			}
			_, e = s.conn.WriteToUDP(b[0:n1], addr)
			if !CheckError(e) {
				return
			}
		}
	}(addr)
	_, e := outgoing.Write(buf)
	CheckError(e)
}
