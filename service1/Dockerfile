FROM golang:1.16.4-alpine3.13
RUN mkdir /build
WORKDIR /build
RUN export GO111MODULE=on
RUN cd /build
RUN go mod init service1
ADD . ./
RUN go build ./service1.go
CMD [ "./service1" ]