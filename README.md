# Nextcloud Talk Bot Framework

A framework to realize Nextcloud Talk chatbots in a client/server model, where sessions exchange data via gRPC stubs.

A Nextcloud Talk Jitsi Bot has been re-written, taking advantage of this framework. It implements the client side and
can be found at [pojntfx/nextcloud-talk-bot-jitsi](https://github.com/pojntfx/nextcloud-talk-bot-jitsi).
The server part is available as `nctalkproxyd`. The source-code is provided in this repo.

Take a look at the following introduction video:

[![thumbnail](https://i3.ytimg.com/vi/WRYlHDGApZo/maxresdefault.jpg)](https://www.youtube.com/watch?v=WRYlHDGApZo)

## Overview

The Nextcloud Talk Bot Framework discribes a client/server infrastructure to realize Nextcloud chatbots,
that interact via gRPC sessions.

* Server side

  `nctalkproxyd` implements a server instance written in the **go** language.

* Client side

  In order to create a chatbot, a client counterpart has to be implemented in any gRPC supported language.
This Client will interacts with `nctalkproxyd` sending and recieving messages. The latter will will take care
of all the heavy lifting (eg. handling the Nextcloud Talk API, keeping track of participants).
`nextcloud-talk-bot-jitsi` is a reference implementation written in **javascript**.

## Installation

### Go Package

A Go package [is available](https://pkg.go.dev/mod/github.com/pojntfx/nextcloud-talk-bot-framework).

### Docker Image

A Docker image is available at [Docker Hub](https://hub.docker.com/r/pojntfx/nctalkproxyd).

### Others

If you're interested in using alternatives like OCI images, see [OCI](./OCI.md).

## Usage

The API will asure fast and secure messsage exchange via gRPC using protocol buffers. The protocol description
itself is defined in [pkg/protos/nextcloud_talk.proto](./pkg/protos/nextcloud_talk.proto).

[`nextcloud-talk-bot-jitsi`](https://github.com/pojntfx/nextcloud-talk-bot-jitsi) is a pretty advanced chatbot implementation,
using this framework. Take it as a reference.

`nctalkproxyd` will integrate itself in the Nextcloud Talk infrastructure while authenticating as a dedicated user.
In order to use the bot, this user (e.g. name it "jitsibot") needs to be added as a participent in every Nextcloud Talk room.
You will handle that as an admin user from within the Nextcloud GUI.

The following code will interconnect a `nctalkproxyd` docker container with a `nextcloud-talk-bot-jitsi`container.
Please adapt variables to meet your production/testing needs. The given values are just examples:

```bash
% docker volume create nctalkproxyd
% docker network create nctalkchatbots
% docker run \
	-p 1969:1969 \
	-v nctalkproxyd:/var/lib/nctalkproxyd \
	-e NCTALKPROXYD_DBPATH=/var/lib/nctalkproxyd \
	-e NCTALKPROXYD_USERNAME=botusername \
	-e NCTALKPROXYD_PASSWORD=botpassword \
	-e NCTALKPROXYD_ADDRREMOTE=https://mynextcloud.com \
	--network nctalkchatbots \
	--name nctalkproxyd \
	-d pojntfx/nctalkproxyd
% docker run \
	-e BOT_JITSI_ADDR=meet.jit.si \
	-e BOT_JITSI_BOT_NAME=botusername \
	-e BOT_JITSI_SLEEP_TIME=20 \
	-e BOT_NCTALKPROXYD_ADDR=nctalkproxyd:1969 \
	-e BOT_JITSI_ROOM_PASSWORD_BYTE_LENGTH=1 \
	-e BOT_COMMANDS=\#videochat,\#videocall,\#custom \
	--network nctalkchatbots \
	-d pojntfx/nextcloud-talk-bot-jitsi
```

## License

Nextcloud Talk Bot Framework (c) 2020 Felix Pojtinger

SPDX-License-Identifier: AGPL-3.0
