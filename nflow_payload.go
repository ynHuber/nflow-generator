package main

import (
	"bytes"
	"encoding/binary"
	"math/rand"
	"net"
	"time"
)

// Start time for this instance, used to compute sysUptime
var StartTime = time.Now().UnixNano()

// current sysUptime in msec - recalculated in CreateNFlowHeader()
var sysUptime uint32 = 0

// Counter of flow packets that have been sent
var flowSequence uint32 = 0

const (
	ICMP_PORT       = 0
	FTP_PORT        = 21
	SSH_PORT        = 22
	DNS_PORT        = 53
	HTTP_PORT       = 80
	HTTPS_PORT      = 443
	NTP_PORT        = 123
	SNMP_PORT       = 161
	IMAPS_PORT      = 993
	MYSQL_PORT      = 3306
	HTTPS_ALT_PORT  = 8080
	P2P_PORT        = 6681
	BITTORRENT_PORT = 6682
	UINT16_MAX      = 65535
	PAYLOAD_AVG_MD  = 1024
	PAYLOAD_AVG_SM  = 256
)

// struct data from fach
type NetflowHeader struct {
	Version        uint16
	FlowCount      uint16
	SysUptime      uint32
	UnixSec        uint32
	UnixMsec       uint32
	FlowSequence   uint32
	EngineType     uint8
	EngineId       uint8
	SampleInterval uint16
}

type NetflowPayload struct {
	SrcIP          uint32
	DstIP          uint32
	NextHopIP      uint32
	SnmpInIndex    uint16
	SnmpOutIndex   uint16
	NumPackets     uint32
	NumOctets      uint32
	SysUptimeStart uint32
	SysUptimeEnd   uint32
	SrcPort        uint16
	DstPort        uint16
	Padding1       uint8
	TcpFlags       uint8
	IpProtocol     uint8
	IpTos          uint8
	SrcAsNumber    uint16
	DstAsNumber    uint16
	SrcPrefixMask  uint8
	DstPrefixMask  uint8
	Padding2       uint16
}

// Complete netflow records
type Netflow struct {
	Header  NetflowHeader
	Records []NetflowPayload
}

// Marshall NetflowData into a buffer
func BuildNFlowPayload(data Netflow) bytes.Buffer {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, &data.Header)
	if err != nil {
		log.Println("Writing netflow header failed:", err)
	}
	for _, record := range data.Records {
		err := binary.Write(buffer, binary.BigEndian, &record)
		if err != nil {
			log.Println("Writing netflow record failed:", err)
		}
	}
	return *buffer
}

// Generate a netflow packet w/ user-defined record count
func GenerateNetflow(bytesPerFlow int, nrOfPackets int, flowDuration time.Duration, trafficDefinition TrafficDefinition) Netflow {
	data := new(Netflow)
	header := CreateNFlowHeader()
	payload := new(NetflowPayload)
	FillCommonFields(payload, uint32(nrOfPackets), uint32(bytesPerFlow), ProtocollForTrafficType(trafficDefinition), rand.Intn(32))

	payload.SrcIP = IPtoUint32("127.0.0.2")
	payload.DstIP = IPtoUint32("127.0.0.1")
	payload.NextHopIP = IPtoUint32("127.0.0.3")

	uptime := int(sysUptime)
	payload.SysUptimeEnd = uint32(uptime)
	payload.SysUptimeStart = payload.SysUptimeEnd - uint32(flowDuration.Milliseconds())

	payload.SrcPrefixMask = 0
	payload.DstPrefixMask = 0

	payload.SrcAsNumber = 0
	payload.DstAsNumber = 553

	payload.SrcPort = uint16(40)
	payload.DstPort = uint16(PortForTrafficType(trafficDefinition))

	records := make([]NetflowPayload, 1)
	records[0] = *payload
	data.Header = header
	data.Records = records
	return *data
}

func PortForTrafficType(trafficDefinition TrafficDefinition) uint16 {
	switch trafficDefinition {
	case FTP:
		return FTP_PORT
	case SSH:
		return SSH_PORT
	case DNS:
		return DNS_PORT
	case HTTP:
		return HTTP_PORT
	case HTTPS:
		return HTTPS_PORT
	case NTP:
		return NTP_PORT
	case SNMP:
		return SNMP_PORT
	case ICMP:
		return ICMP_PORT
	case IMAPS:
		return IMAPS_PORT
	case MYSQL:
		return MYSQL_PORT
	case P2P:
		return P2P_PORT
	case BITTORRENT:
		return BITTORRENT_PORT
	}
	log.Fatal("Unknown traffic definition")
	return 0
}

func ProtocollForTrafficType(trafficDefinition TrafficDefinition) int {
	switch trafficDefinition {
	case FTP:
		return 6
	case SSH:
		return 6
	case DNS:
		return 17
	case HTTP:
		return 6
	case HTTPS:
		return 6
	case NTP:
		return 17
	case SNMP:
		return 17
	case ICMP:
		return 1
	case IMAPS:
		return 6
	case MYSQL:
		return 6
	case P2P:
		return 17
	case BITTORRENT:
		return 17
	default:
	}
	log.Fatal("Unknown traffic definition")
	return 0
}

// patch up the common fields of the packets
func FillCommonFields(
	payload *NetflowPayload,
	numPacket uint32,
	numBytes uint32,
	ipProtocol int,
	srcPrefixMask int) NetflowPayload {

	payload.NumPackets = numPacket
	payload.NumOctets = numBytes

	payload.Padding1 = 0
	payload.IpProtocol = uint8(ipProtocol)
	payload.IpTos = 0
	payload.SrcPrefixMask = uint8(srcPrefixMask)
	payload.DstPrefixMask = uint8(rand.Intn(32))
	payload.Padding2 = 0

	// now handle computed values
	if payload.SrcIP > payload.DstIP { // false-index
		payload.SnmpInIndex = 1
		payload.SnmpOutIndex = 2
	} else {
		payload.SnmpInIndex = 2
		payload.SnmpOutIndex = 1
	}
	return *payload
}

// Generate and initialize netflow header
func CreateNFlowHeader() NetflowHeader {

	t := time.Now().UnixNano()
	sec := t / int64(time.Second)
	nsec := t - sec*int64(time.Second)
	sysUptime = uint32((t-StartTime)/int64(time.Millisecond)) + 1000
	flowSequence++

	h := new(NetflowHeader)
	h.Version = 5
	h.FlowCount = uint16(1)
	h.SysUptime = sysUptime
	h.UnixSec = uint32(sec)
	h.UnixMsec = uint32(nsec)
	h.FlowSequence = flowSequence
	h.EngineType = 1
	h.EngineId = 0
	h.SampleInterval = 1
	return *h
}

func IPtoUint32(s string) uint32 {
	ip := net.ParseIP(s)
	return binary.BigEndian.Uint32(ip.To4())
}
