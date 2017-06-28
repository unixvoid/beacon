Beacon
=======
[![Build Status (Travis)](https://travis-ci.org/unixvoid/beacon.svg?branch=feature%2Ftests)](https://travis-ci.org/unixvoid/beacon)  
Beacon is an API for storing public key-value pairs. It is available for public
use at [beacon.unixvoid.com](https://beacon.unixvoid.com). Check out the
documentation down below to find out how you can register some ID's!

Beacon was designed as a key value discovery service for nsproxy but with global use. Beacon allows
for host discovery by uploading a hosts ip to the beacon server. Once an ip is registered in beacon,
clients are able to pull the ip globally. The beacon api is available for public use globally
at https://beacon.unixvoid.com. This service exposes an api where clients can request a beacon
id and their clients can fetch the data.

Documentation
=============
All documentation is in the [wiki](https://unixvoid.github.io/beacon_docs/)
* [API](https://unixvoid.github.io/beacon_docs/api/)
* [Configuration](https://unixvoid.github.io/beacon_docs/configuration/)
* [Redis](https://unixvoid.github.io/beacon_docs/redis_data_structures/)
* [SSL](https://unixvoid.github.io/beacon_docs/ssl_config/)

Quickstart
==========
To quickly get beacon up and running check out our page on [dockerhub](https://hub.docker.com/r/unixvoid/beacon/)
Or make sure you have [Golang](https://golang.org) and make installed, and use the following make commands:  
* `make deps` to pull down all the 'go gets'
* `make run` to run beacon!
