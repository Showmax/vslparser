#!/bin/sh
set -eu

# The purpose of this script is to build the docker image and to run it as an
# interactive shell with `testdata` directory (the dir this script it located
# in) mounted in the container as /output directory.
#
# As the docker image is supposed to simplify collecting varnishlog varnishlog
# outputs for test purposes, we need to mount local testdata directory into the
# container. That's the main purpose of this script.

cd -- "$(dirname -- "$0")"

readonly image_name="${1:varnish}"

docker build -t "$image_name" .
docker run --rm -it -v "$(pwd)":/output "$image_name" /bin/bash
