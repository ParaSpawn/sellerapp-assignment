FROM golang:latest
ADD ca-certificates.crt /etc/ssl/certs/
RUN mkdir /build
WORKDIR /build
RUN export GO111MODULE=on
RUN cd /build
RUN go mod init service1
RUN go get github.com/ParaSpawn/sellerapp-assignment/service1
RUN go build github.com/ParaSpawn/sellerapp-assignment/service1

EXPOSE 8000

ENTRYPOINT [ "./service1" ]