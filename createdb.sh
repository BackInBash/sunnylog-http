#!/bin/bash
#

docker run -d -e INFLUXDB_DB=sunnylog -p 127.0.0.1:8086:8086 \
      -v $PWD:/var/lib/influxdb \
      influxdb
