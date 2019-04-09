package main

import (
	"encoding/json"
	"github.com/getlantern/systray"
	"io/ioutil"
	"net"
	"os"
)

type Endpoint struct {
	Name        string            `json:"name"`
	IP          string            `json:"IP"`
	Url         string            `json:"url"`
	Fingerprint string            `json:"fingerprint"`
	MEndpoint   *systray.MenuItem `json:"-"`
}

type Config struct {
	Endpoints      []*Endpoint `json:"endpoints"`
	ChosenEndpoint *Endpoint   `json:"-"`
}

type Server struct {
	conn   *net.UDPConn
	config *Config
}

func main() {
	//TODO: register as a service such that it starts on OS login
	systray.Run(onReady, onExit)
}

func Listener(s *Server) { //listen to incoming packets
	Notification("display notification \"Successfully started!\" with title \"DoH\"") //display a notification when successful
	buf := make([]byte, 1024)
	for { //infinite loop
		//TODO: check DNS settings of device from time to time to verify if it still points to 127.0.0.1
		n, addr, err := s.conn.ReadFromUDP(buf)
		if !CheckError(err) {
			continue
		}
		go ProcessRequest(addr, buf[0:n], s) //process request async
	}
}

func getConfig() *Config {
	//TODO: Protect config file, check hash?
	//TODO: add Google (8.8.8.8,8.8.4.4), Quad9 filtered: (9.9.9.9, 149.112.112.112), Quad9 (9.9.9.10, 149.112.112.10)
	var config Config
	configBytes, err := ioutil.ReadFile("config.json")
	if !CheckError(err) {
		os.Exit(1)
	}
	err = json.Unmarshal(configBytes, &config)
	if !CheckError(err) {
		os.Exit(1)
	}
	return &config
}

func createServer() Server { //create server
	serverAddress, err := net.ResolveUDPAddr("udp", ":53") //create address object for server to bind to
	if !CheckError(err) {
		os.Exit(1)
	}
	conn, err := net.ListenUDP("udp", serverAddress)
	if !CheckError(err) {
		os.Exit(1)
	}

	return Server{conn, getConfig()}
}

func onReady() {
	CheckPermissions()   //Checking permissions of runtime, upgrades if necessary
	s := createServer()  //create server
	defer s.conn.Close() //close server if program exits

	systray.SetTitle("DoH")
	systray.SetTooltip("DNS over HTTPS")
	for i, endpoint := range s.config.Endpoints {
		mEndpoint := systray.AddMenuItem(endpoint.Name, endpoint.IP)
		endpoint.MEndpoint = mEndpoint
		go func(s *Server, endpoint *Endpoint) {
			for {
				<-endpoint.MEndpoint.ClickedCh
				s.config.ChosenEndpoint = endpoint
				endpoint.MEndpoint.Check()
				uncheckOthers(s.config.Endpoints, endpoint)
			}
		}(&s, endpoint)
		if i == 0 {
			endpoint.MEndpoint.Check()
			s.config.ChosenEndpoint = endpoint
		}
	}

	mQuit := systray.AddMenuItem("Quit", "Quit")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	Listener(&s) //start listening for infinite time
}

func uncheckOthers(endpoints []*Endpoint, endpoint *Endpoint) {
	for _, otherEndpoint := range endpoints {
		if otherEndpoint != endpoint {
			otherEndpoint.MEndpoint.Uncheck()
		}
	}
}

func onExit() {
	// clean up here
}
