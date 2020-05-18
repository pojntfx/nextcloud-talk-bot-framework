FROM golang:alpine AS build

WORKDIR /app

RUN apk add -u protobuf git
RUN go get github.com/golang/protobuf/protoc-gen-go

COPY . .

RUN go build -o /tmp/nctalkproxyd ./cmd/nctalkproxyd/main.go

FROM alpine

COPY --from=build /tmp/nctalkproxyd /usr/local/bin

CMD /usr/local/bin/nctalkproxyd
