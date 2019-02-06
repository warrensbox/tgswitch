EXE  := tgswitch
PKG  := github.com/warrensbox/tgswitch
VER := $(shell git ls-remote --tags git://github.com/warrensbox/tgswitch | awk '{print $$2}'| awk -F"/" '{print $$3}' | sort -n -t. -k1,1 -k2,2 -k3,3 | tail -n 1)
PATH := build:$(PATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CLIENT_ID := $(CLIENT_ID)
CLIENT_SECRET := $(CLIENT_SECRET)

$(EXE): Gopkg.lock *.go lib/*.go
	go build -v -ldflags "-X main.version=$(VER) -X main.CLIENT_ID=$(CLIENT_ID) -X main.CLIENT_SECRET=$(CLIENT_SECRET)" -o $@ $(PKG)

Gopkg.lock: Gopkg.toml
	dep ensure

.PHONY: release
release: $(EXE) darwin linux

.PHONY: darwin linux 
darwin linux:
	GOOS=$@ go build -ldflags "-X main.version=$(VER) -X main.CLIENT_ID=$(CLIENT_ID) -X main.CLIENT_SECRET=$(CLIENT_SECRET)" -o $(EXE)-$(VER)-$@-$(GOARCH) $(PKG)

.PHONY: clean
clean:
	rm -f $(EXE) $(EXE)-*-*-*

.PHONY: test
test: $(EXE)
	mv $(EXE) build
	go test -v ./...


.PHONEY: dep
dep:
	dep ensure
