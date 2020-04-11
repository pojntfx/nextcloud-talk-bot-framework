# OCI

An Open Containers Initiative (OCI) image for a legacy version of the Nextcloud Talk Jitsi Bot is available at
[Quay.io](https://quay.io/rzerres/nextcloud-talk-jitsi-bot).

Pull the image:

```bash
$ podman pull quay.io/rzerres/nextcloud-talk-jitsi-bot:latest
```

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
