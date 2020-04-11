# Contributing

## Compile binary

On your build system you need an up to date go compiler. Create an executable like this

```bash
$ go build -o /tmp/nextcloud-jitsi-bot main.go
```

## Create an OCI image

The source code does provide a bash-script to create an OCI format based image.
OCI image do allow installation as rootless package

```bash
$ bash -x ./build_container.sh
```

## Uploading to Quay.io

Get the container-id to be uploaded:

```bash
$ imageid=$(buildah image | grep $botname | awk -F " " '{ print $3 }')
```

Tag the container to an image:

```bash
$ buildah commit $imageid quay.io/username/$botname:latest
```

Now login to Quay.io:

```bash
% podman login quay.io
```

Finally upload:

```bash
$ buildah push quay.io/username/$botname:latest
```
