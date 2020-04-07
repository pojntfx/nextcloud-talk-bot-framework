<!-- nextcloud-talk-jitzsi-bot README.md -->
<!-- version: 0.1.1 -->

# Nextcloud Talk Jitsi Bot

A chat bot for Nextcloud Talk that creates [Jitsi](https://jitsi.org) video chat meetings.
It also includes a framework for building Nextcloud Talk chatbots.

## Overview

Nextcloud Talk Jitsi Bot creates ad-hoc Jitsi meetings from Nextcloud Talk sessions.

The talk sessions needs to be aware of the bot. Just create a new `Jitsi-Bot user` inside your trusted `nextcloud` domain.
Its name needs to reflact the `BOT_NEXTCLOUD_USERNAME` environment variable. The bot binary
will resolve it and interact via the `Nextcloud API` to the `Talk` sessions.

When you start a talk session, please add the `Jitsi-Bot user` as a participant of your session.
Once done, you can type `#videocall` or `#videochat`. A newly created link to a Jitsi meeting will be advertised!

## Installation

Precompiled binaries are provided from the following locations.

### Go module

A Go package is available from githubs module repro [pkg.go.dev](https://pkg.go.dev/mod/github.com/pojntfx/nextcloud-talk-jitsi-bot).

### Docker Image

A Docker image is available at [Docker Hub](https://hub.docker.com/r/pojntfx/nextcloud-talk-jitsi-bot).

Pull the image

	docker pull hub.docker.com/pojntfx/nextcloud-talk-jitsi-bot

### OCI Image

An Open Containers Initiative (OCI) image is available at
[Quay.IO](https://quay.io/rzerres/nextcloud-talk-jitsi-bot).

Pull the image

	podman pull quay.io/rzerres/nextcloud-talk-jitsi-bot:latest

<!--
#### Howto uploading to Quay.io

Get the container-id to be uploaded

	imageid=$(buildah image | grep $botname | awk -F " " '{ print $3 }')

Tag the container to an image

	buildah commit $imageid quay.io/username/$botname:latest

Now login to Quay.io

	podman login quay.io
	username: <username>
	password: <password>

Finally upload

	buildah push quay.io/username/$botname:latest

-->


## Build

### Compile binary

On your build system you need an up to date go compiler. Create an executable like this

	go build -o /tmp/nextcloud-jitsi-bot main.go


### create a docker image

The source code does provide a dockerfile, that will handle the creation of a docker image

### create an OCI image

The source code does provide a bash-script to create an OCI format based image.
OCI image do allow installation as rootless package

	bash -x ./build_container.sh

## Usage

You might also add a group
`nextcloud-talk-jitsi-bot` gets all needed parameters via customisable environment variables.
Please adapt them as appropriate and start the bot binary with reference to this environment.

### Docker image

Please adapt then environment varialbles as needed. You can run the image like:

```bash
% docker volume create nextcloud-talk-jitsi-bot # This is required so that messages don't get send twice
% docker run \
	-e BOT_NEXTCLOUD_USERNAME="jitsibot" \
	-e BOT_NEXTCLOUD_PASSWORD="password" \
	-e BOT_DB_LOCATION="/run/nextcloud-jitsi-bot" \
	-e BOT_NEXTCLOUD_URL="https://localhost:8443" \
	-e BOT_JITSI_URL="https://meet.jit.si" \
	-v nextcloud-talk-jitsi-bot:/var/lib/nextcloud-jitsi-bot \
	-d \
	pojntfx/nextcloud-talk-jitsi-bot
```

### OCI image

On a `systemd` capable distro, you OCI images can be managed combining a `systemd.service` with the `podman`
binary. If you are not familiar with `podman` yet, you might simply put: `alias docker=podman`. Beside using
it as a drop-in replacement for docker, there are a couple of advanteges:

* OCI compliant
* rootless and root mode
* daemonless
* direct interaction with Container Registy, Containers, Imag Storage and runc

To create the systemd.service, run

	podman create --detach --name nextcloud-talk-jitsi-bot nextcloud-talk-jitsi-bot:latest
	podman generate systemd --name nextcloud-talk-jitsi-bot > /etc/systemd/system/nextcloud-talk-jitsi-bot.service

Have a look at [running containers with podman]( https://www.redhat.com/sysadmin/podman-shareable-systemd-services)
to get more insightdetails.

	systemctl edit --full nextcloud-talk-jitsi-bot.service

Adapt the environment variables with appropriate values. They are configured with following defaults:

	BOT_NEXTCLOUD_USERNAME="jitsibot"
	BOT_NEXTCLOUD_PASSWORD="password"
	BOT_DB_LOCATION="/run/nextcloud-jitsi-bot"
	BOT_NEXTCLOUD_URL="https://localhost:8443"
	BOT_JITSI_URL="https://meet.jit.si"

Finally start the service.

	systemctl enable nextcloud-talk-jitsi-bot.service
	systemctl start nextcloud-talk-jitsi-bot.service

## License

Nextcloud Talk Jitsi Bot (c) 2020 Felicitas Pojtinger

SPDX-License-Identifier: AGPL-3.0
