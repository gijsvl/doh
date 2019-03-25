package main

import (
	"encoding/json"
	"github.com/getlantern/systray"
	"io/ioutil"
	"net"
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
	systray.Run(onReady, onExit)
}

func Listener(s *Server) { //listen to incoming packets
	Notification("display notification \"Successfully started!\" with title \"DoH\"") //display a notification when successful
	buf := make([]byte, 1024)
	for { //infinite loop
		//TODO: some calls fail... QUIC Protocol, maybe? google.com
		// https://tools.ietf.org/id/draft-huitema-quic-dnsoquic-03.html
		// https://datatracker.ietf.org/meeting/99/materials/slides-99-dprive-dns-over-quic-01
		n, addr, err := s.conn.ReadFromUDP(buf)
		CheckError(err)
		go ProcessRequest(addr, buf[0:n], s) //process request async
	}
}

func getConfig() *Config {
	//TODO: Protect config file, check hash?
	var config Config
	configBytes, err := ioutil.ReadFile("config.json")
	CheckError(err)
	err = json.Unmarshal(configBytes, &config)
	CheckError(err)
	return &config
}

func createServer() Server { //create server
	serverAddress, err := net.ResolveUDPAddr("udp", ":53") //create address object for server to bind to
	CheckError(err)
	conn, err := net.ListenUDP("udp", serverAddress)
	CheckError(err)

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
