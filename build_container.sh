#!/bin/sh

set -x

###
# buildah: creating a container for Nextcloud-Talk proxy daemon (OCI format)
###

push_image=0
mount_image=1

# variables
author='Felix Pojtinger @pojntfx'
maintainer='Ralf Zerres @rzerres'
botname=nctalkproxyd
buildroot=/run/buildah
prefix=/usr/local

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

# get dependencies
buildah run $mycontainer apk add --no-cache --virtual .build-deps protobuf git
buildah run $mycontainer go get github.com/golang/protobuf/protoc-gen-go
buildah copy $mycontainer .

# create the binary
buildah run $mycontainer go build -o $buildroot/$botname ./cmd/$botname/main.go

# prepare the destination container
buildah run $mycontainer mkdir -p $prefix/etc
buildah run $mycontainer cp $buildroot/$botname $prefix/bin/$botname
buildah run $mycontainer cp $buildroot/examples/$botname.yaml $prefix/etc/$botname.yaml
buildah run $mycontainer ln -s $prefix/etc/$botname.yaml /etc/$botname.yaml
buildah run $mycontainer mkdir /run/$botname

# allow manual adaptions
if [ "$mount_image" -eq 1 ]; then
read -p "Mount the created container for interactive adaptions (y/n)? " -t 15 doMount
	# for rootless mode: run in user namespace
	if [ "$doMount" = "y" ]; then
		echo "You are inside the container tree. Build environment hasn't been flushed yet!"
		if [ $(id -u) -eq 0 ]; then
			# execution in root mode
			mountpoint=$(buildah mount $mycontainer)
			cd $mountpoint
			du -sh *
			sh
			buildah umount $mycontainer
		else
			# execution in rootless mode
			# not starting with usernamespace, default to isolate the filesystem with chroot
			#ENV _BUILDAH_STARTED_IN_USERNS="" BUILDAH_ISOLATION=chroot
			buildah unshare --mount containerID du -sh ${containerID}/*
			#buildah unshare du -sh $mountpoint/*
			buildah unshare sh
			#buildah unshare umount $mountpoint
		fi
	fi
fi

# cleanup dependencies
buildah run $mycontainer apk del .build-deps

# cleanup build environment
buildah config --workingdir="/" $mycontainer
buildah run $mycontainer rm -rf $buildroot
buildah run $mycontainer rm -rf go
buildah run $mycontainer rm -rf usr/local/go usr/local/lib usr/local/share root/.cache

# tag the container to an image name, sign it. on success remove container
imageid=$(buildah commit \
	--rm \
	--squash \
	$mycontainer $repo/$username/$botname)
	#--sign-by $signkey \

# tag our new image with an alternate name
#buildah tag $botname nctjb

if [ "$push_image" -eq 1 ]; then
	# push image to Quay Container Registry
	#imageid=$(buildah images | grep $repo/$username/$botname | awk -F ' ' '{print $3}')
	buildah push \
		--authfile $authfile \
		$imageid docker://$repo/$username/$botname
fi
