GOC=go build
GOFLAGS=-a -ldflags '-s'
CGOR=CGO_ENABLED=0
VER_NUM=latest
DOCKER_OPTIONS="--no-cache"
IMAGE_NAME=mfaltys/beacon:$(VER_NUM)
REDIS_DB_HOST_DIR="/tmp/"

all: beacon
beacon: beacon.go
	$(GOC) beacon.go
run: beacon.go
	go run beacon.go
stat: beacon.go
	$(CGOR) $(GOC) $(GOFLAGS) beacon.go
install: stat
	cp beacon /usr/bin

docker:
	$(MAKE) stat
	mkdir stage.tmp/
	cp beacon stage.tmp/
	cp deps/rootfs.tar.gz stage.tmp/
	cp deps/Dockerfile stage.tmp/
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
	rm -f beacon
	rm -rf stage.tmp/
#CGO_ENABLED=0 go build -a -ldflags '-s' beacon.go
