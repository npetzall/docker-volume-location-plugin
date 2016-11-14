APP									= docker-volume-location-plugin

VERSION 						?= HEAD
BUILD               ?= $(shell git rev-parse HEAD)

PLATFORMS           = linux_amd64 linux_arm darwin_amd64

FLAGS_all = CGO_ENABLED=0

FLAGS_linux_amd64   = $(FLAGS_all) GOOS=linux GOARCH=amd64
FLAGS_linux_arm     = $(FLAGS_all) GOOS=linux GOARCH=arm
FLAGS_darwin_amd64  = $(FLAGS_all) GOOS=darwin GOARCH=amd64


msg=@printf "\n\033[0;01m>>> %s\033[0m\n" $1

.DEFAULT_GOAL := build

get-deps:
	go get github.com/tools/godep
	godep restore
.PHONY: deps

build: get-deps
	$(call msg,"Build binary")
	$(FLAGS_all) go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}" -o bin/${APP} .
	./bin/${APP} -version
.PHONY: build

test: get-deps
	$(call msg,"Run tests")
	go test -v .
.PHONY: test

clean:
	$(call msg,"Clean directory")
	rm -rf bin
	rm -rf dist
.PHONY: clean

build-all: get-deps $(foreach PLATFORM,$(PLATFORMS),dist/$(PLATFORM)/$(APP))
.PHONY: build-all

dist: notHEAD build-all \
$(foreach PLATFORM,$(PLATFORMS),dist/$(APP)-$(VERSION)-$(PLATFORM).tar.gz)
.PHONY:	dist

dist/%/$(APP):
	$(call msg,"Build binary for $*")
	rm -f $@
	mkdir -p $(dir $@)
	$(FLAGS_$*) go build -ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}" -o dist/$*/${APP} .

dist/$(APP)-$(VERSION)-%.tar.gz: notHEAD
	$(call msg,"Create TAR for $*")
	rm -f $@
	mkdir -p $(dir $@)
	tar czf $@ -C dist/$* .

notHEAD:
	@ if [ "$(VERSION)" = "HEAD" ]; then \
	 		echo "Not allowed with VERSION == HEAD"; \
			exit 1; \
		fi
