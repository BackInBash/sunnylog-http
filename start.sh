#!/bin/sh
while true
do
    /opt/sunnylog-v2 --solarUrl $SUNNYBOYURL --solarPassword $SUNNYBOYPASSWD --influxUrl $INFLUXAPI --influxToken $INFLUXTOKEN
    # Sleep 5 min
    sleep 300
done