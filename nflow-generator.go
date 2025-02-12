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
	"time"

	"github.com/jessevdk/go-flags"
)

type Proto int

const (
	FTP Proto = iota + 1
	SSH
	DNS
	HTTP
	HTTPS
	NTP
	SNMP
	IMAPS
	MYSQL
	HTTPS_ALT
	P2P
	BITTORRENT
)

var opts struct {
	CollectorIP        string `short:"t" long:"target" description:"target ip address of the netflow collector"`
	CollectorPort      string `short:"p" long:"port" description:"port number of the target netflow collector"`
	FalseIndex         bool   `short:"f" long:"false-index" description:"generate false SNMP interface indexes, otherwise set to 0"`
	Help               bool   `short:"h" long:"help" description:"show nflow-generator help"`
	ReportingIntervall string `short:"i" long:"reporting-intervall" description:"intervall in which the flow should be send to the collector"`
	NrOfPackets        int    `short:"n" long:"packets-per-flow" description:"set the number of reported packets per flow"`
	BytesPerFlow       int    `short:"b" long:"bytes-per-flow" description:"set the number of reported bytes per flow"`
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
		showUsage()
		os.Exit(1)
	}
	if opts.BytesPerFlow == 0 {
		showUsage()
		os.Exit(1)
	}
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
		data := GenerateNetflow(opts.BytesPerFlow, opts.NrOfPackets, intervall)
		buffer := BuildNFlowPayload(data)
		_, err = conn.Write(buffer.Bytes())
		if err != nil {
			log.Fatal("Error connecting to the target collector: ", err)
		}
		time.Sleep(intervall - time.Since(start))
		start = flowEnd
	}
}

func showUsage() {
	var usage string
	usage = `
Usage:
  main [OPTIONS] [collector IP address] [collector port number]

  Send mock Netflow version 5 data to designated collector IP & port.
  Time stamps in all datagrams are set to UTC.

Application Options:
  -t, --target= target ip address of the netflow collector
  -p, --port=   port number of the target netflow collector
  -s, --spike run a second thread generating a spike for the specified protocol
    protocol options are as follows:
        ftp - generates tcp/21
        ssh  - generates tcp/22
        dns - generates udp/54
        http - generates tcp/80
        https - generates tcp/443
        ntp - generates udp/123
        snmp - generates ufp/161
        imaps - generates tcp/993
        mysql - generates tcp/3306
        https_alt - generates tcp/8080
        p2p - generates udp/6681
        bittorrent - generates udp/6682
  -f, --false-index generate a false snmp index values of 1 or 2. The default is 0. (Optional)
  -c, --flow-count set the number of flows to generate in each iteration. The default is 16. (Optional)

Example Usage:

    -first build from source (one time)
    go build   

    -generate default flows to device 172.16.86.138, port 9995
    ./nflow-generator -t 172.16.86.138 -p 9995 

    -generate default flows along with a spike in the specified protocol:
    ./nflow-generator -t 172.16.86.138 -p 9995 -s ssh

    -generate default flows with "false index" settings for snmp interfaces 
    ./nflow-generator -t 172.16.86.138 -p 9995 -f

    -generate default flows with up to 256 flows
    ./nflow-generator -c 128 -t 172.16.86.138 -p 9995

Help Options:
  -h, --help    Show this help message
  `
	fmt.Print(usage)
}
