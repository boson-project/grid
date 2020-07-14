FROM golang:alpine as build
RUN apk add make git gcc g++ 
COPY . /src
WORKDIR /src
RUN make

# Base Image
FROM alpine:latest
RUN apk --update add supervisor ca-certificates
COPY --from=build /src/faas /bin/
COPY ./supervisor.conf /
EXPOSE 8080

CMD ["supervisord","-c","/supervisor.conf"]

# Attempt 1
#  Supposedly has systemd enabled
# FROM registry.access.redhat.com/ubi8-init
# RUN dnf -y install httpd; dnf clean all; systemctl enable httpd
# CMD [ "/sbin/init" ]

# Attempt 2:
# # Setup:
#  mount -t cgroup2 none /sys/fs/cgroup
#  pacman -S crun
#  podman build -t systemd .
# # Dockerfile:
# FROM fedora
# RUN dnf -y install httpd; dnf clean all; systemctl enable httpd
# EXPOSE 1111
# CMD [ "/sbin/init" ]




