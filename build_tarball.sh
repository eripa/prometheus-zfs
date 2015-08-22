#!/bin/bash
#
# Simple bash script for building binaries for all relevant platforms

SCRIPT_DIR=$(dirname $0)
cd ${SCRIPT_DIR}

# Build
declare -a TARGETS=(darwin linux solaris freebsd)
for target in ${TARGETS[@]} ; do
  output="prometheus-zfs-${target}"
  echo "Building for ${target}, output bin/${output}"
  GOOS=${target}
  GOARCH=amd64
  go build -o bin/${output}
done

# Create a tar-ball for release
DIR_NAME=${PWD##*/} # name of current directory, presumably prometheus-zfs
VERSION=$(git describe --abbrev=0 --tags 2> /dev/null)
if [ "$?" -ne 0 ] ; then
    # No tag, use commit hash
    HASH=$(git rev-parse HEAD)
    VERSION=${HASH:0:7}
fi

cd ../
TARBALL="prometheus-zfs-${VERSION}.tar.gz"
tar -cf ${TARBALL} --exclude=.git -vz ${DIR_NAME}
echo "Created: ${PWD}/${TARBALL}"
