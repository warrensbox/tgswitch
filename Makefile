EXE  := tgswitch
PKG  := github.com/warrensbox/tgswitch
VER := $(shell { git ls-remote --tags . 2>/dev/null || git ls-remote --tags git@github.com:warrensbox/tgswitch.git; } | awk '{if ($$2 ~ "\\^\\{\\}$$") next; print vers[split($$2,vers,"\\/")]}' | sort -n -t. -k1,1 -k2,2 -k3,3 | tail -1)
PATH := build:$(PATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

$(EXE): go.mod *.go lib/*.go
	go build -v -ldflags "-X main.version=$(VER)" -o $@ $(PKG)

.PHONY: release
release: $(EXE) darwin linux

.PHONY: darwin linux
darwin linux:
	GOOS=$@ go build -ldflags "-X main.version=$(VER)" -o $(EXE)-$(VER)-$@-$(GOARCH) $(PKG)

.PHONY: clean
clean:
	rm -f $(EXE) $(EXE)-*-*-*

.PHONY: test
test: $(EXE)
	mkdir -p build
	mv $(EXE) build
	go test -v ./...

.PHONY: install
install: $(EXE)
	mkdir -p ~/bin
	mv $(EXE) ~/bin

.PHONY: docs
docs:
	cd docs; bundle install --path vendor/bundler; bundle exec jekyll build -c _config.yml; cd ..

.PHONY: version
version:
	@echo $(VER)

