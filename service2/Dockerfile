FROM golang:1.16.4-alpine3.13
RUN mkdir /build
WORKDIR /build
RUN export GO111MODULE=on
RUN cd /build
RUN go mod init service2
ADD . ./
RUN go build ./service2.go
CMD [ "./service2" ]