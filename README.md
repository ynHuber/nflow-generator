# nflow-generator

Netflow generator that produces fixed netflow mock data to test monitoring systems

Based on https://github.com/nerdalert/nflow-generator

#### Usecase
This generator allows the creation of specific netflow data to monitor netflow monitoring platform.

The use of configuration files allows sending multiple specific netflows to multiple netflow recievers

#### Example configuration:
```
generators:
  - collector:
      ip: 127.0.0.1
      port: 3055
    flows:
      - traffic-type: ntp
        reporting-intervall: 1s
        packets-per-flow: 10
        bytes-per-flow: 280
        sample-intervall: 5
      - traffic-type: dns
        reporting-intervall: 1s
        packets-per-flow: 10
        bytes-per-flow: 280
        sample-intervall: 5
  - collector:
      ip: second.netflow.receiver.domain
      port: 3055
    flows:
      - traffic-type: icmp
        reporting-intervall: 1s
        packets-per-flow: 10
        bytes-per-flow: 10
        sample-intervall: 1
      - traffic-type: CLDAP
        reporting-intervall: 1s
        packets-per-flow: 10
        bytes-per-flow: 280
        sample-intervall: 5
```
## Integrated SNMP Mock Server
A mock SNMP server can be set up using a configuration file. 
This might be using using the nflowgenerator for testing a flow monitoring system that obtains additional data via snmp.
The SNMP Mock Server can configure static return values for multiple OIDs:

```
snmp:
  ip: localhost
  port: 161
  community: community
  mockOIDs:
    - oid: ".1.3.6.1.2.1.31.1.1.1.1.1"
      value: "interface1"
      type: "string"
    - oid: ".1.3.6.1.2.1.31.1.1.1.15.1"
      value: "100000000000"
      type: "uint64"
    - oid: ".1.3.6.1.2.1.31.1.1.1.18.1"
      value: "interface1 - 100 Gbps"
      type: "string"
    - oid: ".1.3.6.1.2.1.31.1.1.1.1.2"
      value: "interface2"
      type: "string"
    - oid: ".1.3.6.1.2.1.31.1.1.1.15.2"
      value: "100000000000"
      type: "uint64"
    - oid: ".1.3.6.1.2.1.31.1.1.1.18.2"
      value: "interface2 - 100 Gbps"
      type: "string"
```

Atm only string and unit64 are implemented as types