# The pirpose of this Dockerfile is to simplify generation of varnishlog
# outputs by simply wrapping the whole varnish setup into a docker container.

FROM archlinux:latest
RUN pacman -Syu --noconfirm && pacman -S pacman-contrib --noconfirm && paccache -r
RUN pacman -S varnish tmux vim openbsd-netcat --noconfirm && paccache -r
RUN echo -e '#!/bin/bash\n\n/usr/sbin/varnishd -a :6081 -f /etc/varnish/default.vcl -s malloc,256m' > run_varnish.sh && chmod +x run_varnish.sh
