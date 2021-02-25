FROM golang:alpine3.13

COPY build/sunnylog /opt/sunnylog
COPY start.sh /opt/start.sh

ENTRYPOINT [ "/opt/start.sh" ]