FROM golang:alpine AS build

WORKDIR /app

COPY . .

RUN go build -o /tmp/nextcloud-talk-jitsi-bot main.go

FROM alpine

COPY --from=build /tmp/nextcloud-talk-jitsi-bot /usr/local/bin

CMD /usr/local/bin/nextcloud-talk-jitsi-bot
