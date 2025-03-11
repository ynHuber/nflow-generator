# nflow-generator

Netflow generator that produces fixed netflow mock data to test monitoring systems

Based on https://github.com/nerdalert/nflow-generator

#### Usecase
This generator allows the creation of specific netflow data to monitor netflow monitoring platform.

The use of configuration files allows sending multiple specific netflows to multiple netflow recievers

#### Example configuration:
```
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
