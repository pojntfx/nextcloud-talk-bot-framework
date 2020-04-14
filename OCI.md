# OCI

An Open Containers Initiative (OCI) image that provides `nxtalkproxyd` is available at
[Quay.io](https://quay.io/rzerres/nctalkproxyd).

The images can be combined with the bot image of `nctalk-bot-jitsi`. Its up to you, to
choose an adequate pod mode.

* umbrella pod: a single pod, where both images are legal citizens
* multi pod:  each image is executed its own pod instance (eg. to address scalability, cluster-awareness, etc)

Please talk into account, that internetworking between pods is ony supported for images running in root mode
(as of podman <= v1.8.2)

## Installation

Pull the image:

```bash
$ podman pull quay.io/rzerres/nctalk-proxyd:latest
```

Check the network setup

```bash
$ podman network ls
```bash

if you need a new definition, go ahead and create one with your prefered driver (default: "bridge")

```bash
$ podman network create  --driver bridge nctalkproxyd
```

## Systemd handling (multi pod)

On a `systemd` capable distro, you OCI images can be managed combining a `systemd.service` with the `podman`
binary. If you are not familiar with `podman` yet, you might simply put: `alias docker=podman`. Beside using
it as a drop-in replacement for docker, there are a couple of advantages:

- OCI compliant
- rootless and root mode
- daemonless
- direct interaction with Container Registy, Containers, Image Storage and runc

To create the systemd.service, run

```bash
$ podman create --detach --name nctalkproxyd nctalkproxy:latest -u <username> -p <password> -r "https://your.nextcloud.url"
$ podman generate systemd --name nctalkproxyd > /etc/systemd/system/nctalkproxyd.service
```

Have a look at [running containers with podman](https://www.redhat.com/sysadmin/podman-shareable-systemd-services)
to get more insightdetails.

## Adaptation

#The required parameters for `nctalkproxyd` can be adapted either via a config file, via config
parameters or via corresponding environment variables. The latter take precedence.

```bash
$ systemctl edit --full nctalk-proxyd.service
```

The config file is preset with the following defaults:

```bash
$ cat /etc/nctalkproxyd.yaml

nctalkproxyd:
  addrLocale: :1969
  addrRemote: https://mynextcloud.com
  username: botusername
  password: botpassword
  dbpath: /var/lib/nctalkproxyd
```

Adapt the environment variables with appropriate values.

```env
NCTALKPROXYD_DBPATH=/var/lib/nctalkproxyd
NCTALKPROXYD_USERNAME=botusername
NCTALKPROXYD_PASSWORD=botpassword
NCTALKPROXYD_ADDRLOCAL=:1969
NCTALKPROXYD_ADDRREMOTE=https://mynextcloud.com
```

Finally start the service.

```bash
$ systemctl enable nctalkproxy.service
$ systemctl start nctalkproxy.service
```
