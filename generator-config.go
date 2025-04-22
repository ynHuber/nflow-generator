package main

import "gopkg.in/yaml.v2"

type OIDType string

const (
	String OIDType = "string"
	Uint64 OIDType = "uint64"
)

type Configuration struct {
	SNMPServerConfiguration *SNMPServerConfigruation  `yaml:"snmp,omitempty"`
	GeneratorConfigurations []*GeneratorConfiguration `yaml:"generators"`
}

type GeneratorConfiguration struct {
	CollectorConfiguration *CollectorConfiguration `yaml:"collector"`
	FlowConfigurations     []*FlowConfiguration    `yaml:"flows"`
}

type CollectorConfiguration struct {
	CollectorIP   string `yaml:"ip"`
	CollectorPort string `yaml:"port"`
}

type FlowConfiguration struct {
	ReportingIntervall string `yaml:"reporting-intervall"`
	NrOfPackets        int    `yaml:"packets-per-flow"`
	BytesPerFlow       int    `yaml:"bytes-per-flow"`
	TrafficType        string `yaml:"traffic-type"`
	SampleInterval     uint16 `yaml:"sample-intervall" default:"1"`
}

type SNMPServerConfigruation struct {
	CollectorIP   string     `yaml:"ip"  default:""`
	CollectorPort string     `yaml:"port" default:"161"`
	Community     string     `yaml:"community"`
	MockOIDs      []*MockOID `yaml:"mockOIDs"`
	//Additional Parameters for SNMPv3
	V3Username                 string `yaml:"username" default:"v3Username"`
	V3AuthenticationPassphrase string `yaml:"auth_pass" default:"v3AuthenticationPassphrase"`
	V3PrivacyPassphrase        string `yaml:"privacy_pass" default:"v3PrivacyPassphrase"`
}

type MockOID struct {
	OID   string  `yaml:"oid"`
	Value string  `yaml:"value"`
	Type  OIDType `yaml:"type"`
}

func ParseConfiguration(b []byte) Configuration {
	configuration := new(Configuration)
	err := yaml.Unmarshal(b, configuration)
	if err != nil {
		log.Fatalf("[error] Error parsing configuration YAML: %v", err)
	}
	return *configuration
}
