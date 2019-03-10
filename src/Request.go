package main

import (
	"encoding/base64"
	"io/ioutil"
	"net"
)

func ProcessRequest(addr *net.UDPAddr, buf []byte, ServerConn *net.UDPConn) {
	resolvedQuery := DNSOverHTTPSRequest(base64.StdEncoding.EncodeToString(buf))      //request DNS record over HTTPS
	if resolvedQuery != nil {
		_, err := ServerConn.WriteToUDP(resolvedQuery, addr)
		CheckError(err)
	}
}

func DNSOverHTTPSRequest(record string) []byte {
	//QUERY OVER HTTPS
	//TODO: parameterize the https-endpoint
	client := NewClient([]byte{61, 149, 205, 222, 84, 64, 203, 239, 45, 4, 169, 54, 59, 30, 133, 238, 50, 37, 159, 58, 246, 99, 57, 176, 169, 205, 201, 159, 39, 215, 160, 44})
	res, err := client.Get("https://1.1.1.1/dns-query?dns=" + record)
	CheckError(err)
	body, err := ioutil.ReadAll(res.Body)
	CheckError(err)

	return body
}
