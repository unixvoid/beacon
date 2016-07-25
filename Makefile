GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0
VER_NUM=latest
DOCKER_OPTIONS="--no-cache"
IMAGE_NAME=docker.io/unixvoid/beacon:$(VER_NUM)
REDIS_DB_HOST_DIR="/tmp/"

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
	$(CGOR) $(GOC) $(GOFLAGS) -o bin/beacon beacon/*.go

test:
	go test -v beacon/*.go

install: stat
	cp beacon /usr/bin

docker:
	$(MAKE) stat
	mkdir stage.tmp/
	cp beacon stage.tmp/
	cp deps/rootfs.tar.gz stage.tmp/
	cp deps/Dockerfile stage.tmp/
	sed -i "s/<DIFF>/$(shell git rev-parse HEAD | head -c 10)/g" stage.tmp/Dockerfile
	chmod +x deps/run.sh
	cp deps/run.sh stage.tmp/
	cp config.gcfg stage.tmp/
	cd stage.tmp/ && \
		sudo docker build $(DOCKER_OPTIONS) -t $(IMAGE_NAME) .
	@echo "$(IMAGE_NAME) built"
dockerrun:
	sudo docker run \
		-d \
		-p 8808:8808 \
		--name beacon \
		-v $(REDIS_DB_HOST_DIR):/redisbackup/:rw \
		mfaltys/beacon
	sudo docker logs -f beacon

clean:
	rm -rf bin/
	rm -rf stage.tmp/
#CGO_ENABLED=0 go build -a -ldflags '-s' beacon.go
