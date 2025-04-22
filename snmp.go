package main

import (
	"strconv"
	"sync"

	"github.com/gosnmp/gosnmp"
	"github.com/rs/zerolog/log"
	"github.com/slayercat/GoSNMPServer"
)

func runMockSNMPServer(conf *SNMPServerConfigruation, startUpWait *sync.WaitGroup) {
	oids := []*GoSNMPServer.PDUValueControlItem{}

OIDLooP:
	for _, oid := range conf.MockOIDs {

		var msg interface{}
		var oidType gosnmp.Asn1BER

		switch oid.Type {
		case String:
			oidType = gosnmp.OctetString
			msg = GoSNMPServer.Asn1OctetStringWrap(oid.Value)
		case Uint64:
			oidType = gosnmp.Counter64
			oidValueUint, err := strconv.ParseUint(oid.Value, 10, 0)
			if err != nil {
				log.Warn().Err(err)
				continue OIDLooP
			}
			msg = GoSNMPServer.Asn1Counter64Unwrap(oidValueUint)
		default:
			log.Warn().Msgf("unknown oid Type %s", oid.Type)
			continue OIDLooP
		}

		oidObj := GoSNMPServer.PDUValueControlItem{
			OID:      oid.OID,
			Type:     oidType,
			OnGet:    func() (value interface{}, err error) { return msg, nil },
			Document: "ifIndex",
		}
		oids = append(oids, &oidObj)

	}
	var logger = &ZeroLogger{}

	master := GoSNMPServer.MasterAgent{
		Logger: logger,
		SecurityConfig: GoSNMPServer.SecurityConfig{
			AuthoritativeEngineBoots: 1,
			Users: []gosnmp.UsmSecurityParameters{
				{
					UserName:                 conf.V3Username,
					AuthenticationProtocol:   gosnmp.MD5,
					PrivacyProtocol:          gosnmp.DES,
					AuthenticationPassphrase: conf.V3AuthenticationPassphrase,
					PrivacyPassphrase:        conf.V3PrivacyPassphrase,
				},
			},
		},
		SubAgents: []*GoSNMPServer.SubAgent{
			{
				CommunityIDs: []string{conf.Community},
				OIDs:         oids,
			},
		},
	}
	server := GoSNMPServer.NewSNMPServer(master)
	serverAddress := conf.CollectorIP + ":" + conf.CollectorPort
	err := server.ListenUDP("udp", serverAddress)
	if err != nil {
		log.Error().Err(err).Msg("Error in listen")
	}
	go server.ServeForever()
	startUpWait.Done()
}
