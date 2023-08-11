FROM golang:alpine AS builder

WORKDIR /build

COPY . . 

WORKDIR /build/cmd 

RUN go build -o test main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/cmd/test /build/test

ARG file_name

COPY $file_name .

ENV file_name=$file_name

ENTRYPOINT ./test $file_name
