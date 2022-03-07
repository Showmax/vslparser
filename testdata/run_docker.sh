#!/bin/sh

# The purpose of this script is to build the docker image and to run it as an
# interactive shell with `testdata` directory (the dir this script it located
# in) mounted in the container as /output directory.

readonly image_name="${1:varnish}"

docker build -t "$image_name" .
docker run --rm -it -v "$(pwd)":/output "$image_name" /bin/bash
