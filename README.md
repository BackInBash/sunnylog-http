# sunnylog-http
Download production data from SMA Sunny Boy using the HTTP interface.

**Tested against SUNNY BOY 1.5 with Firmware 3.10.7.R**

## Usage:

```bash
./sunnylog-http --solarUrl http://[SunnyBoyUrl] --solarPassword [SunnyBoy User Password] --influxUrl http://[InfluxDB HTTP API]
```

### Supported Operating Systems:
+ Windows x64
+ Linux x64
+ Linux ARM

## Build

To Build the Project the following InfluxDB Client Branch must be used `master-1.x`.

1. To Switch the Branch `go get` all Dependencies
1. Change to InfluxDB Client Folder `cd ~/go/src/github.com/influxdata/influxdb`
1. Switch Git Branch `git checkout master-1.x`
1. Build the Project `go build` 💥
