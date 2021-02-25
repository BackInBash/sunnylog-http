#!/bin/sh
while true
do
    /opt/sunnylog --solarUrl $SUNNYBOYURL --solarPassword $SUNNYBOYPASSWD --influxUrl $INFLUXAPI --influxUser $INFLUXUSER --influxPass $INFLUXPASS
    # Sleep 5 min
    sleep 300
done