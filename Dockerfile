FROM golang:alpine3.13

COPY sunnylog.go /opt/sunnylog.go
COPY go.mod /opt/go.mod
COPY go.sum /opt/go.sum
COPY start.sh /opt/start.sh

RUN su -c "cd /opt && go get"

RUN su -c "cd /opt && go build && chmod +x sunnylog-v2"

ENTRYPOINT [ "/opt/start.sh" ]