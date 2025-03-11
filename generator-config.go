package main

import "gopkg.in/yaml.v2"

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

func ParseGeneratorConfigurations(b []byte) []GeneratorConfiguration {
	generatorConfigurations := new([]GeneratorConfiguration)
	err := yaml.Unmarshal(b, generatorConfigurations)
	if err != nil {
		log.Fatalf("[error] Error parsing configuration YAML: %v", err)
	}
	return *generatorConfigurations
}
