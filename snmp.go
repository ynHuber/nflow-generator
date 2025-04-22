package main

import (
	"os"
	"strconv"
	"sync"

	"github.com/gosnmp/gosnmp"
	"github.com/sirupsen/logrus"
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
				log.Warn(err)
				continue OIDLooP
			}
			log.Info(oidValueUint)
			msg = GoSNMPServer.Asn1Counter64Unwrap(oidValueUint)
		default:
			log.Warn("unknown oid Type " + oid.Type)
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
	var log = logrus.New()
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel

	master := GoSNMPServer.MasterAgent{
		Logger: log,
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
		log.Printf("Error in listen: %+v", err)
	}
	go server.ServeForever()
	startUpWait.Done()
}
