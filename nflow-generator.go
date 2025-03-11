// Run using:
// go run nflow-generator.go nflow_logging.go nflow_payload.go  -t 172.16.86.138 -p 9995
// Or:
// go build
// ./nflow-generator -t <ip> -p <port>
package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

type TrafficDefinition int

const (
	FTP TrafficDefinition = iota + 1
	SSH
	DNS
	HTTP
	HTTPS
	NTP
	SNMP
	ICMP
	IMAPS
	MYSQL
	P2P
	BITTORRENT
)

var opts struct {
	CollectorIP        string `short:"t" long:"target" description:"target ip address of the netflow collector"`
	CollectorPort      string `short:"p" long:"port" description:"port number of the target netflow collector"`
	Help               bool   `short:"h" long:"help" description:"show nflow-generator help"`
	ReportingIntervall string `short:"i" long:"reporting-intervall" description:"intervall in which the flow should be send to the collector"`
	NrOfPackets        int    `short:"n" long:"packets-per-flow" description:"set the number of reported packets per flow"`
	BytesPerFlow       int    `short:"b" long:"bytes-per-flow" description:"set the number of reported bytes per flow"`
	TrafficDef         string `short:"d" long:"defiendTraffic" description:"generated netflow traffic definition (default=ntp)"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		showUsage()
		os.Exit(1)
	}
	if opts.Help == true {
		showUsage()
		os.Exit(1)
	}
	if opts.CollectorIP == "" || opts.CollectorPort == "" {
		showUsage()
		os.Exit(1)
	}
	if opts.ReportingIntervall == "" {
		showUsage()
		os.Exit(1)
	}
	intervall, err := time.ParseDuration(opts.ReportingIntervall)
	if err != nil {
		log.Fatal("Failed to parse reporting intervall \""+opts.ReportingIntervall+"\" due to ", err.Error())
	}

	if opts.NrOfPackets == 0 {
		opts.NrOfPackets = 10
	}
	if opts.BytesPerFlow == 0 {
		showUsage()
		os.Exit(1)
	}

	trafficDef := parseTrafficDefinition(opts.TrafficDef)

	collector := opts.CollectorIP + ":" + opts.CollectorPort
	udpAddr, err := net.ResolveUDPAddr("udp", collector)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal("Error connecting to the target collector: ", err)
	}
	log.Infof("sending netflow data to a collector ip: %s and port: %s. \n"+
		"Use ctrl^c to terminate the app.", opts.CollectorIP, opts.CollectorPort)

	start := time.Now()
	for {
		flowEnd := start.Add(intervall)
		data := GenerateNetflow(opts.BytesPerFlow, opts.NrOfPackets, intervall, trafficDef)
		buffer := BuildNFlowPayload(data)
		_, err = conn.Write(buffer.Bytes())
		if err != nil {
			log.Fatal("Error connecting to the target collector: ", err)
		}
		time.Sleep(intervall - time.Since(start))
		start = flowEnd
	}
}

func parseTrafficDefinition(trafficDefinitionString string) TrafficDefinition {
	var trafficDef TrafficDefinition

	switch strings.ToUpper(trafficDefinitionString) {
	case "":
		//default
		trafficDef = NTP
	case "FTP":
		trafficDef = FTP
	case "SSH":
		trafficDef = SSH
	case "DNS":
		trafficDef = DNS
	case "HTTP":
		trafficDef = HTTP
	case "HTTPS":
		trafficDef = HTTPS
	case "NTP":
		trafficDef = NTP
	case "SNMP":
		trafficDef = SNMP
	case "ICMP":
		trafficDef = ICMP
	case "IMAPS":
		trafficDef = IMAPS
	case "MYSQL":
		trafficDef = MYSQL
	case "P2P":
		trafficDef = P2P
	case "BITTORRENT":
		trafficDef = BITTORRENT

	default:
		log.Fatal("Failed to parse netflow traffic definition" + opts.TrafficDef)
		showUsage()
		os.Exit(1)
	}
	return trafficDef
}

func showUsage() {
	var usage string
	usage = `
Usage:
  main [OPTIONS] [collector IP address] [collector port number]

  Send mock Netflow version 5 data to designated collector IP & port.
  Time stamps in all datagrams are set to UTC.

Application Options:
  -b --bytes-per-flow= number of reported bytes per flow
  -d --netflow-traffic-definition= generated netflow traffic definition (default=ntp)
    protocol options are as follows:
        ftp - generates tcp/21
        ssh  - generates tcp/22
        dns - generates udp/54
        http - generates tcp/80
        https - generates tcp/443
        ntp - generates udp/123
        snmp - generates ufp/161
        icmp - generates icmp
        imaps - generates tcp/993
        mysql - generates tcp/3306
        p2p - generates udp/6681
        bittorrent - generates udp/6682
  -n, --nr-of-packets-per-flow= number of packets per flow. (default=10)
  -p, --port= port number of the target netflow collector
  -r  --reporting-intervall= intervall in which flow messages should be send to the collector
  -t, --target= target ip address of the netflow collector

Help Options:
  -h, --help    Show this help message
  `
	fmt.Print(usage)
}
