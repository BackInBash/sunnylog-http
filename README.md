# sunnylog-http
Download production data from SMA Sunny Boy using the HTTP interface.

**Tested against SUNNY BOY 1.5 with Firmware 3.10.7.R**

## Usage:

### Standalone Binary
```bash
./sunnylog-http --solarUrl http://[SunnyBoyUrl] --solarPassword [SunnyBoy User Password] --influxUrl http://[InfluxDB HTTP API]
```

### Docker Container

```bash
docker run -d -e SUNNYBOYURL=http://[SunnyBoyUrl] \
 -e SUNNYBOYPASSWD=[SunnyBoy User Password] \
 -e INFLUXAPI=http://[InfluxDB HTTP API] \
 -e INFLUXTOKEN=[InfluxDB User Token] \
 ghcr.io/backinbash/sunnylog-http/sunnylog:v2
```

### Supported Operating Systems:
+ Windows x64
+ Linux x64
+ Linux ARM

## Build Binary

1. Clone the Repo
1. `go get`
1. Build the Project `go build` 💥

## Build Container
1. Clone the Repo
1. `docker build -t sunnylog .`