# The purpose of this Dockerfile is to simplify generation of varnishlog
# outputs for testing purposed. To avoid the need to install Varnish during
# development process, we have decided to rather separate this complexity into
# a Docker image.
#
# Expected usage of this Dockerfile is to build & run it locally using an
# interactive shell (see ./run_docker.sh script in this directory). Once you
# run inside the container, you can use the /run_varnish.sh to start Varnish on
# port :6081. Then you can run varnishlog using whichever configuration (set of
# options) you want to generate test for and use curl to generate logs desired
# logs.

FROM archlinux:latest
RUN pacman -Sy archlinux-keyring --noconfirm && pacman -Syu --noconfirm && pacman -S pacman-contrib --noconfirm && paccache -rk0
RUN pacman -S varnish tmux vim openbsd-netcat --noconfirm && paccache -rk0
RUN echo -e '#!/bin/bash\n\n/usr/sbin/varnishd -a :6081 -f /etc/varnish/default.vcl -s malloc,256m' > run_varnish.sh && chmod +x run_varnish.sh
