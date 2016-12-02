GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0
VER_NUM=latest
DOCKER_OPTIONS="--no-cache"
IMAGE_NAME=docker.io/unixvoid/beacon:$(VER_NUM)
REDIS_DB_HOST_DIR="/tmp/"
GIT_HASH=$(shell git rev-parse HEAD | head -c 10)

all: beacon

beacon: beacon.go
	$(GOC) beacon.go

dependencies:
	go get github.com/gorilla/mux
	go get golang.org/x/crypto/sha3
	go get gopkg.in/gcfg.v1
	go get gopkg.in/redis.v3
	go get github.com/unixvoid/glogger

run:
	cd beacon && go run \
	beacon.go \
	provision.go \
	remove.go \
	rotate.go \
	update.go

stat:
	mkdir -p bin/
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/beacon-$(GIT_HASH)-linux-amd64 beacon/*.go

test:
	go test -v beacon/*.go

install: stat
	cp beacon /usr/bin

docker:
	$(MAKE) stat
	mkdir stage.tmp/
	cp bin/beacon* stage.tmp/
	cp deps/rootfs.tar.gz stage.tmp/
	cp deps/Dockerfile stage.tmp/
	sed -i "s/<DIFF>/$(GIT_HASH)/g" stage.tmp/Dockerfile
	chmod +x deps/run.sh
	cp deps/run.sh stage.tmp/
	cp beacon/config.gcfg stage.tmp/
	cd stage.tmp/ && \
		sudo docker build $(DOCKER_OPTIONS) -t $(IMAGE_NAME) .
	@echo "$(IMAGE_NAME) built"

dockerrun:
	sudo docker run \
		-d \
		-p 8808:8808 \
		--name beacon \
		-v $(REDIS_DB_HOST_DIR):/redisbackup/:rw \
		unixvoid/beacon
	sudo docker logs -f beacon

aci:
	$(MAKE) stat
	mkdir -p stage.tmp/beacon-layout/rootfs/
	tar -zxf deps/rootfs.tar.gz -C stage.tmp/beacon-layout/rootfs/
	cp bin/beacon* stage.tmp/beacon-layout/rootfs/beacon
	chmod +x deps/run.sh
	cp deps/run.sh stage.tmp/beacon-layout/rootfs/
	sed -i "s/<DIFF>/$(GIT_HASH)/g" stage.tmp/beacon-layout/rootfs/run.sh
	cp beacon/config.gcfg stage.tmp/beacon-layout/rootfs/
	cp deps/manifest.json stage.tmp/beacon-layout/manifest
	cd stage.tmp/ && \
		actool build beacon-layout beacon.aci && \
		mv beacon.aci ../
	@echo "beacon.aci built"

travisaci:
	wget https://github.com/appc/spec/releases/download/v0.8.7/appc-v0.8.7.tar.gz
	tar -zxf appc-v0.8.7.tar.gz
	$(MAKE) stat
	mkdir -p stage.tmp/beacon-layout/rootfs/
	tar -zxf deps/rootfs.tar.gz -C stage.tmp/beacon-layout/rootfs/
	cp bin/beacon* stage.tmp/beacon-layout/rootfs/beacon
	chmod +x deps/run.sh
	cp deps/run.sh stage.tmp/beacon-layout/rootfs/
	sed -i "s/<DIFF>/$(GIT_HASH)/g" stage.tmp/beacon-layout/rootfs/run.sh
	cp beacon/config.gcfg stage.tmp/beacon-layout/rootfs/
	cp deps/manifest.json stage.tmp/beacon-layout/manifest
	cd stage.tmp/ && \
		../appc-v0.8.7/actool build beacon-layout beacon.aci && \
		mv beacon.aci ../
	@echo "beacon.aci built"

clean:
	rm -rf bin/
	rm -rf stage.tmp/
#CGO_ENABLED=0 go build -a -ldflags '-s' beacon.go
