#!/bin/bash

set -x

###
# buildah: creating a Nextcloud-Talk jitsi bot container (OCI format)
###

push_image=false

# variables
author='Felicitas Pojtinger @pojntfx'
maintainer='Ralf Zerres @rzerres'
botname=nextcloud-talk-jitsi-bot
buildroot=/run/buildah
prefixbin=/usr/local/bin

# signing with GPG Package-Key
signkey=1EC4BE4FF2A6C9F4DDDF30F33C5F485DBD250D66

# authenticate at quay.io with robot account (env: REGISTRY_AUTH_FILE)
repo=quay.io
username=rzerres
authfile=$HOME/.docker/config.json

# create a container -> docker image golang (flavor: alpine)
mycontainer=$(buildah from --name $botname golang:alpine)

# adapt containers metadata
buildah config --author="$author" $mycontainer
buildah config --label name=$botname $mycontainer
buildah config --label maintainer="$maintainer" $mycontainer
buildah config --env BOT_NEXTCLOUD_USERNAME="jitsibot" $mycontainer
buildah config --env BOT_NEXTCLOUD_PASSWORD="password" $mycontainer
buildah config --env BOT_DB_LOCATION="/run/$botname" $mycontainer
buildah config --env BOT_NEXTCLOUD_URL="https://localhost" $mycontainer
buildah config --env BOT_JITSI_URL="https://meet.jit.si" $mycontainer

# create the build environment inside the container
#buildah config --env GOPATH="$buildroot" $mycontainer
buildah config --entrypoint "[ \"/usr/local/bin/$botname\" ]" $mycontainer
buildah config --workingdir="$buildroot" $mycontainer

# create the binary
#buildah run $mycontainer mkdir -p $buildroot/pkg
buildah copy $mycontainer src
buildah run $mycontainer go build -o $buildroot/$botname main.go

# prepare the destination container
buildah run $mycontainer cp $buildroot/$botname $prefixbin/$botname
buildah run $mycontainer mkdir /run/$botname

# cleanup build environment
buildah config --workingdir="/" $mycontainer
buildah run $mycontainer rm -rf $buildroot
buildah run $mycontainer rm -rf go
buildah run $mycontainer rm -rf usr/local/go usr/local/lib usr/local/share

# not starting with usernamespace, default to isolate the filesystem with chroot
#ENV _BUILDAH_STARTED_IN_USERNS="" BUILDAH_ISOLATION=chroot

# manual testings
# for rootless mode: run in user namespace
#buildah unshare

#mountpoint=$(buildah mount $mycontainer)
#ls -l $mountpoint/$prefixbin
#buildah umount $mycontainer

# tag the container to an image name, sign it. on success remove container
imageid=$(buildah commit \
	--rm \
	--squash \
	$mycontainer $repo/$username/$botname)
	#--sign-by $signkey \

# tag our new image with an alternate name
#buildah tag $botname nctjb

if [ "$push_image" == "true" ]; then
	# push image to Quay Container Registry
	#imageid=$(buildah images | grep $repo/$username/$botname | awk -F ' ' '{print $3}')
	buildah push \
		--authfile $authfile \
		$imageid docker://$repo/$username/$botname
fi
