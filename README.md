# nextcloud-talk-jitsi-bot

A chat bot for Nextcloud Talk that creates [Jitsi](https://jitsi.org) video chat meetings.

## Overview

Nextcloud Talk Jitsi Bot creates ad-hoc Jitsi meetings from Nextcloud Talk. Just add the bot to the conversation, type `#videocall` or `#videochat` and get a link to a Jitsi meeting!

## Installation

A Docker image is available on [Docker Hub](https://hub.docker.com/r/pojntfx/nextcloud-talk-jitsi-bot).

## Usage

First, create an account for the bot in Nextcloud.

Configure the bot with the following env variables:

```bash
NCTB_SERVER="https://nx6978.your-storageshare.
de"
NCTB_USER="jitsi--bot"
NCTB_PASS="secretpassword"
NCTB_JITSI_URL="https://meet.jit.si"
```

And run it:

```bash
docker run -e NCTB_SERVER=$NCTB_SERVER -e NCTB_USER=$NCTB_USER -e NCTB_PASS=$NCTB_PASS -e NCTB_JITSI_URL=$NCTB_JITSI_URL pojntfx/nextcloud-talk-jitsi-bot
```

## License

nextcloud-talk-jitsi-bot (c) 2020 Felicitas Pojtinger

SPDX-License-Identifier: AGPL-3.0
