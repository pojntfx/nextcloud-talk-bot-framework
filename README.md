# Nextcloud Talk Jitsi Bot

A chat bot for Nextcloud Talk that creates [Jitsi](https://jitsi.org) video chat meetings. It also includes a framework for building Nextcloud Talk chatbots.

## Overview

Nextcloud Talk Jitsi Bot creates ad-hoc Jitsi meetings from Nextcloud Talk. Just add the bot to the conversation, type `#videocall` or `#videochat` and get a link to a Jitsi meeting!

## Installation

### Go Package

A Go package [is available](https://pkg.go.dev/github.com/pojntfx/nextcloud-talk-jitsi-bot).

### Docker Image

A Docker image is available on [Docker Hub](https://hub.docker.com/r/pojntfx/nextcloud-talk-jitsi-bot).

## Usage

First, create an account for the bot in Nextcloud.

Then, run it:

```bash
% docker volume create nextcloud-talk-jitsi-bot # This is required so that messages don't get send twice
% docker run \
    -e BOT_NEXTCLOUD_URL="https://nx6978.your-storageshare.de" \
    -e BOT_NEXTCLOUD_USERNAME="botusername" \
    -e BOT_NEXTCLOUD_PASSWORD="botpassword" \
    -e BOT_DB_LOCATION="/var/lib/nextcloud-jitsi-bot" \
    -e BOT_JITSI_URL="https://meet.jit.si" \
    -v nextcloud-talk-jitsi-bot:/var/lib/nextcloud-jitsi-bot \
    -d \
    pojntfx/nextcloud-talk-jitsi-bot
```

## License

Nextcloud Talk Jitsi Bot (c) 2020 Felix Pojtinger

SPDX-License-Identifier: AGPL-3.0
