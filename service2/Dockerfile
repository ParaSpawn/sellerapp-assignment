FROM golang:latest

RUN mkdir /build
WORKDIR /build

RUN export GO111MODULE=on
RUN cd /build
RUN go mod init service2
RUN go get github.com/ParaSpawn/sellerapp-assignment/service2
RUN go build github.com/ParaSpawn/sellerapp-assignment/service2

EXPOSE 8010

ENTRYPOINT [ "./service2" ]