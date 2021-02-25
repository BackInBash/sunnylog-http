FROM golang:alpine3.13

COPY sunnylog.go /opt/sunnylog.go
COPY start.sh /opt/start.sh

RUN su -c "cd /opt && go mod init sunnylog && go mod tidy && go get"

RUN su -c "cd /opt && go build && chmod +x sunnylog"

ENTRYPOINT [ "/opt/start.sh" ]