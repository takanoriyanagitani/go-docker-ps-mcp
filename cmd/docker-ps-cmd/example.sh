#!/bin/sh

pattern='^.*$'

port=12029

dsock=/var/run/docker.sock
#dsock=~/docker.sock

./docker-ps-cmd \
	-port ${port} \
	-container-name-pattern "${pattern}" \
	-docker-host "${dsock}"
