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
	"sync"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type TrafficType int

const (
	FTP TrafficType = iota + 1
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
	CLDAP
)

type CommandLineOptions struct {
	CollectorIP        string `short:"t" long:"target" description:"target ip address of the netflow collector"`
	CollectorPort      string `short:"p" long:"port" description:"port number of the target netflow collector"`
	Config             string `short:"c" long:"config" description:"path to configuration file"`
	Help               bool   `short:"h" long:"help" description:"show nflow-generator help"`
	ReportingIntervall string `short:"i" long:"reporting-intervall" description:"intervall in which the flow should be send to the collector"`
	NrOfPackets        int    `short:"n" long:"packets-per-flow" description:"set the number of reported packets per flow"`
	BytesPerFlow       int    `short:"b" long:"bytes-per-flow" description:"set the number of reported bytes per flow"`
	TrafficType        string `short:"d" long:"defiendTraffic" description:"generated netflow traffic definition" default:"ntp"`
	SampleInterval     uint16 `short:"s" long:"sample-intervall" description:"Sample intervall stated in the netflow headers" default:"1"`
}

func main() {
	opts := CommandLineOptions{}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.DateTime})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	_, err := flags.Parse(&opts)
	if err != nil {
		showUsage()
		os.Exit(1)
	}
	if opts.Help == true {
		showUsage()
		os.Exit(1)
	}

	if opts.Config == "" {
		runSimpleNetflowGeneration(opts)
	} else {
		runWithConfig(opts.Config)
	}

}

func runWithConfig(configFilePath string) {
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Printf("[error] reading config file: %s", err)
		return
	}
	configuration := ParseConfiguration(configFile)

	if configuration.SNMPServerConfiguration != nil {
		snmpStartupWait := &sync.WaitGroup{}
		snmpStartupWait.Add(1)
		go runMockSNMPServer(configuration.SNMPServerConfiguration, snmpStartupWait)
		snmpStartupWait.Wait()
	}

	wg := &sync.WaitGroup{}
	for _, generatorConfiguration := range configuration.GeneratorConfigurations {
		go runGenerator(generatorConfiguration, wg)
	}

	// wait for wait group to register in runGenerator
	time.Sleep(time.Duration(time.Second))
	wg.Wait()
}

func runSimpleNetflowGeneration(opts CommandLineOptions) {
	if opts.CollectorIP == "" || opts.CollectorPort == "" {
		showUsage()
		os.Exit(1)
	}
	if opts.ReportingIntervall == "" {
		showUsage()
		os.Exit(1)
	}

	if opts.NrOfPackets == 0 {
		opts.NrOfPackets = 10
	}
	if opts.BytesPerFlow == 0 {
		showUsage()
		os.Exit(1)
	}

	generatorConfiguration := GeneratorConfiguration{
		CollectorConfiguration: &CollectorConfiguration{
			CollectorIP:   opts.CollectorIP,
			CollectorPort: opts.CollectorPort,
		},
		FlowConfigurations: []*FlowConfiguration{
			{
				ReportingIntervall: opts.ReportingIntervall,
				NrOfPackets:        opts.NrOfPackets,
				BytesPerFlow:       opts.BytesPerFlow,
				TrafficType:        opts.TrafficType,
				SampleInterval:     opts.SampleInterval,
			},
		},
	}

	wg := &sync.WaitGroup{}
	runGenerator(&generatorConfiguration, wg)
	wg.Wait()
}

func runGenerator(generatorConfiguration *GeneratorConfiguration, wg *sync.WaitGroup) {
	collectorConfiguration := generatorConfiguration.CollectorConfiguration
	collector := collectorConfiguration.CollectorIP + ":" + collectorConfiguration.CollectorPort
	udpAddr, err := net.ResolveUDPAddr("udp", collector)
	if err != nil {
		log.Fatal().Err(err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer func() {
		wg.Wait()
		conn.Close()
	}()
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to the target collector:")
	}

	for _, flowDefinition := range generatorConfiguration.FlowConfigurations {
		wg.Add(1)
		go sendFlows(conn, flowDefinition, wg)
	}

	log.Info().Msgf("sending netflow data to a collector ip: %s and port: %s. - "+
		"Use ctrl^c to terminate the app.", collectorConfiguration.CollectorIP, collectorConfiguration.CollectorPort)
}

func sendFlows(conn *net.UDPConn, flowDefinition *FlowConfiguration, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
	}()
	trafficType := parseTrafficType(flowDefinition.TrafficType)

	intervall, err := time.ParseDuration(flowDefinition.ReportingIntervall)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to parse reporting intervall \"%s\" due to ", flowDefinition.ReportingIntervall)
	}

	start := time.Now()
	for {
		flowEnd := start.Add(intervall)
		data := GenerateNetflow(flowDefinition.BytesPerFlow, flowDefinition.NrOfPackets, intervall, trafficType, flowDefinition.SampleInterval)
		buffer := BuildNFlowPayload(data)
		_, err = conn.Write(buffer.Bytes())
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to the target collector")
		}
		time.Sleep(intervall - time.Since(start))
		start = flowEnd
	}
}

func parseTrafficType(trafficTypeString string) TrafficType {
	var trafficType TrafficType

	switch strings.ToUpper(trafficTypeString) {
	case "":
		//default
		trafficType = NTP
	case "FTP":
		trafficType = FTP
	case "SSH":
		trafficType = SSH
	case "DNS":
		trafficType = DNS
	case "HTTP":
		trafficType = HTTP
	case "HTTPS":
		trafficType = HTTPS
	case "NTP":
		trafficType = NTP
	case "SNMP":
		trafficType = SNMP
	case "ICMP":
		trafficType = ICMP
	case "IMAPS":
		trafficType = IMAPS
	case "MYSQL":
		trafficType = MYSQL
	case "P2P":
		trafficType = P2P
	case "BITTORRENT":
		trafficType = BITTORRENT
	case "CLDAP":
		trafficType = CLDAP

	default:
		log.Fatal().Msgf("Failed to parse netflow traffic definition %s", trafficTypeString)
		showUsage()
		os.Exit(1)
	}
	return trafficType
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
        ssh - generates tcp/22
		cldap - generates udp/389
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
  -r, --reporting-intervall= intervall in which flow messages should be send to the collector
  -s, --sample-intervall= Sample intervall stated in the netflow headers. (default=1)
  -t, --target= target ip address of the netflow collector

Using a configuration file
	The option -c, --config can be used to run the generator using a configuration file.
	This allows specifiying multiple flow definitions that are run in parallel

Help Options:
  -h, --help    Show this help message
  `
	fmt.Print(usage)
}
