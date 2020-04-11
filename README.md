# Nextcloud Talk Bot Framework

A framework for writing Nextcloud Talk chatbots with every language that supports gRPC.

Looking for the Nextcloud Talk Jitsi Bot? It has been re-written as a client for this framework at [pojntfx/nextcloud-talk-bot-jitsi](https://github.com/pojntfx/nextcloud-talk-bot-jitsi)!

## Overview

The Nextcloud Talk Bot Framework provides a way to create chatbots for Nextcloud in any language that supports gRPC. To do so, `nxtalkproxyd` - a streaming gRPC API for Nextcloud Talk - is the primary part of the framework; in order for you to create a chatbot, you just have to write a client for `nxtalkproxyd`, which will take care of all the heavy lifting for you!

## Installation

### Go Package

A Go package [is available](https://pkg.go.dev/mod/github.com/pojntfx/nextcloud-talk-bot-framework).

### Docker Image

A Docker image is available at [Docker Hub](https://hub.docker.com/r/pojntfx/nxtalkproxyd).

### Others

If you're interested in using alternatives like OCI images, see [OCI](./OCI.md).

## Usage

As explained above, all you have to to write a chatbot is to implement a client for `nxtalkproxyd`! Take a look at [pkg/protos/nextcloud_talk.proto](./pkg/protos/nextcloud_talk.proto) for the protocol. A pretty advanced chatbot that is based on this framework is the [Nextcloud Talk Jitsi Bot](https://github.com/pojntfx/nextcloud-talk-bot-jitsi), so if you want to have a quick start take a look at the repo.

`nxtalkproxyd` requires a user to work with; every Nextcloud Talk room that should be able to use the bots connected to it has to add this user.

To then start using your bot, you can connect your bot to `nxtalkproxyd` like so (this is the way that the Nextcloud Talk Jitsi Bot does it; don't forget to change i.e. the username and password, it is just an example):

```bash
% docker volume create nxtalkproxyd
% docker network create nxtalkchatbots
% docker run \
    -p 1969:1969 \
    -v nxtalkproxyd:/var/lib/nxtalkproxyd \
    -e NXTALKPROXYD_NXTALKPROXYD_DBPATH=/var/lib/nxtalkproxyd \
    -e NXTALKPROXYD_NXTALKPROXYD_USERNAME=botusername \
    -e NXTALKPROXYD_NXTALKPROXYD_PASSWORD=botpassword \
    -e NXTALKPROXYD_NXTALKPROXYD_RADDR=https://examplenextcloud.com \
    --network nxtalkchatbots \
    --name nxtalkproxyd \
    -d pojntfx/nxtalkproxyd
% docker run \
    -e BOT_JITSI_ADDR=meet.jit.si \
    -e BOT_JITSI_BOT_NAME=botusername \
    -e BOT_JITSI_SLEEP_TIME=20 \
    -e BOT_NXTALKPROXYD_ADDR=nxtalkproxyd:1969 \
    -e BOT_JITSI_ROOM_PASSWORD_BYTE_LENGTH=1 \
    -e BOT_COMMANDS=\#videochat,\#videocall,\#custom \
    --network nxtalkchatbots \
    -d pojntfx/nextcloud-talk-bot-jitsi
```

## License

Nextcloud Talk Bot Framework (c) 2020 Felicitas Pojtinger

SPDX-License-Identifier: AGPL-3.0
