package main

import (
	"net"
)

func main() {
	CheckPermissions()       //Checking permissions of runtime, upgrades if necessary
	ServerConn := Server()   //create server
	defer ServerConn.Close() //close server if program exits
	Listener(ServerConn)     //start listening for infinite time
}

func Listener(conn *net.UDPConn) { //listen to incoming packets
	Notification("display notification \"Successfully started!\" with title \"DoH\"") //display a notification when successful
	buf := make([]byte, 1024)
	for { //infinite loop
		n, addr, err := conn.ReadFromUDP(buf)
		CheckError(err)
		go ProcessRequest(addr, buf[0:n], conn) //process request async
	}
}

func Server() *net.UDPConn { //create server
	serverAddress, err := net.ResolveUDPAddr("udp", ":53") //create address object for server to bind to
	CheckError(err)                                        //check error, if found, log and exit
	conn, err := net.ListenUDP("udp", serverAddress)
	CheckError(err) //check error, if found, log and exit
	return conn
}
