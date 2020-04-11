# Nextcloud Talk Jitsi Bot

A chat bot for Nextcloud Talk that creates [Jitsi](https://jitsi.org) video chat meetings.
It also includes a framework for building Nextcloud Talk chatbots.

## Overview

Nextcloud Talk Jitsi Bot creates ad-hoc Jitsi meetings from Nextcloud Talk sessions.

The talk sessions needs to be aware of the bot. Just create a new user for the bot inside your trusted Nextcloud domain.
Its username needs to reflect the `BOT_NEXTCLOUD_USERNAME` environment variable. The bot binary will resolve it and interact with Nextcloud Talk over it's API.

When you start a talk session, please add the bot user you've added as a participant of your room. Once done, you can type `#videocall` or `#videochat`. A newly created link to a Jitsi meeting will be advertised!

## Installation

### Go module

A Go package is available from the GitHub Go Module Registry: [pkg.go.dev](https://pkg.go.dev/mod/github.com/pojntfx/nextcloud-talk-jitsi-bot).

### Docker Image

A Docker image is available at [Docker Hub](https://hub.docker.com/r/pojntfx/nextcloud-talk-jitsi-bot).

Pull the image:

```bash
$ docker pull hub.docker.com/pojntfx/nextcloud-talk-jitsi-bot
```

### OCI Image

An Open Containers Initiative (OCI) image is available at
[Quay.io](https://quay.io/rzerres/nextcloud-talk-jitsi-bot).

Pull the image:

```bash
$ podman pull quay.io/rzerres/nextcloud-talk-jitsi-bot:latest
```

## Usage

You might also add a group
`nextcloud-talk-jitsi-bot` gets all needed parameters via customisable environment variables.
Please adapt them as appropriate and start the bot binary with reference to this environment.

### Docker image

Please adapt then environment varialbles as needed. You can run the image like:

```bash
$ docker volume create nextcloud-talk-jitsi-bot # This is required so that messages don't get send twice
$ docker run \
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

- OCI compliant
- rootless and root mode
- daemonless
- direct interaction with Container Registy, Containers, Image Storage and runc

To create the systemd.service, run

```bash
$ podman create --detach --name nextcloud-talk-jitsi-bot nextcloud-talk-jitsi-bot:latest
$ podman generate systemd --name nextcloud-talk-jitsi-bot > /etc/systemd/system/nextcloud-talk-jitsi-bot.service
```

Have a look at [running containers with podman](https://www.redhat.com/sysadmin/podman-shareable-systemd-services)
to get more insightdetails.

```bash
$ systemctl edit --full nextcloud-talk-jitsi-bot.service
```

Adapt the environment variables with appropriate values. They are configured with following defaults:

```env
BOT_NEXTCLOUD_USERNAME="jitsibot"
BOT_NEXTCLOUD_PASSWORD="password"
BOT_DB_LOCATION="/run/nextcloud-jitsi-bot"
BOT_NEXTCLOUD_URL="https://localhost:8443"
BOT_JITSI_URL="https://meet.jit.si"
```

Finally start the service.

```bash
$ systemctl enable nextcloud-talk-jitsi-bot.service
$ systemctl start nextcloud-talk-jitsi-bot.service
```

## Build

See [CONTRIBUTING](./CONTRIBUTING.md) for more information.

## License

Nextcloud Talk Jitsi Bot (c) 2020 Felix Pojtinger

SPDX-License-Identifier: AGPL-3.0
