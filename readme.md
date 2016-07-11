Beacon
=======
[![Build Status (Travis)](https://travis-ci.org/unixvoid/beacon.svg?branch=feature%2Ftests)](https://travis-ci.org/unixvoid/beacon)
Beacon is a key-value discovery service designed for nsproxy but with global use. This project is mainly taken from bitnuke, with code from here Beacon allows for host discovery by uploading a hosts ip to the beacon server. Once an ip is registered in beacon, clients are able to pull the ip globally. The beacon api is available for public use globally at https://beacon.unixvoid.com. This service exposes an api where clients can request a beacon id and their clients can fetch the data.

Documentation
=============
All documentation is in the [github wiki](https://github.com/unixvoid/beacon/wiki)
* [Configuration](https://github.com/unixvoid/beacon/wiki/Configuration)
* [API](https://github.com/unixvoid/beacon/wiki/API)
* [Redis](https://github.com/unixvoid/beacon/wiki/Redis)
* [SSL](https://github.com/unixvoid/beacon/wiki/SSL)

Quickstart
==========
To quickly get beacon up and running check out our page on [dockerhub](https://hub.docker.com/r/unixvoid/beacon/)
Or make sure you have [Golang](https://golang.org) and make installed, and use the following make commands:  
* `make deps` to pull down all the 'go gets'
* `make run` to run beacon!
