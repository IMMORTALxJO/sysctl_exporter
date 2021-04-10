# sysctl exporter

That application exports sysctl parameters in prometheus format.
Parameters with numbered values are supported only.

### How to run

Build and run docker image:
```
# docker build . -t sysctl_exporter
# docker run -d --net host --pid host sysctl_exporter
# curl http://localhost:9141/metrics
...
sysctl_parameter{param="vm.overcommit_kbytes"} 0
sysctl_parameter{param="vm.overcommit_memory"} 1
sysctl_parameter{param="vm.overcommit_ratio"} 50
...
```
or use already built image from docker hub
```
docker run -d --net host --pid host immortalxjo/sysctl_exporter:latest
```

### Usage
```
  -exclude string
    	RegExp for skipping sysctl parameters
  -include string
    	RegExp for sysctl parameters (default ".*")
  -listen-address string
    	Address to listen on for telemetry (default ":9141")
  -log-level string
    	Verbosity of logging (default "info")
  -metrics-prefix string
    	Prefix of prometheus metrics (default "sysctl")
```

### Format
There are two variants of metrics
#### one number parameter
`net.ipv4.udp_rmem_min = 4096`
generate metric
`sysctl_parameter{param="net.ipv4.udp_rmem_min"} 4096`
#### multiple numbers parameter
`net.ipv4.tcp_rmem = 4096	87380	6291456`
produce metrics
```
sysctl_parameter{column="0",param="net.ipv4.tcp_rmem"} 4096
sysctl_parameter{column="1",param="net.ipv4.tcp_rmem"} 87380
sysctl_parameter{column="2",param="net.ipv4.tcp_rmem"} 6.291456e+06
```
