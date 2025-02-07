Netflow generator that produces fixed netflow mock data to test monitoring systems
Based on https://github.com/nerdalert/nflow-generator

Create mock Netflow for a 100G DoS using `go run . -t 127.0.0.1 -p 9001 -i "10ms" -n 1000000 -b 1000000000`