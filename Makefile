EXE  := tgswitch
PKG  := github.com/warrensbox/tgswitch
VER  := $(shell { git ls-remote --tags . 2>/dev/null || git ls-remote --tags git@github.com:warrensbox/tgswitch.git; } | awk '{if ($$2 ~ "\\^\\{\\}$$") next; print vers[split($$2,vers,"\\/")]}' | sort -n -t. -k1,1 -k2,2 -k3,3 | tail -1)
PATH := build:$(PATH)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

GOMOD     = on
DIR       = ./build
GOOSX     = darwin
GOOSLINUX = linux
OSXBIN    = $(DIR)/$(EXE)-$(VER)-$(GOOSX)-$(GOARCH)
LINUXBIN  = $(DIR)/$(EXE)-$(VER)-$(GOOSLINUX)-$(GOARCH)

$(EXE): go.mod *.go lib/*.go
	go build -v -ldflags "-X main.version=$(VER)" -o $@ $(PKG)

.PHONY: release
release: $(EXE) darwin linux

.PHONY: $(OSXBIN)
$(OSXBIN):
	GO111MODULE=$(GOMOD) GOOS=$(GOOSX) go build -ldflags "-X main.version=$(VER)" -o $(OSXBIN) $(PKG)

.PHONY: darwin
darwin: $(OSXBIN)
	chmod +x $(OSXBIN)

.PHONY: $(LINUXBIN)
$(LINUXBIN):
	GO111MODULE=$(GOMOD) GOOS=$(GOOSLINUX) go build -ldflags "-X main.version=$(VER)" -o $(LINUXBIN) $(PKG)

.PHONY: linux
linux: $(LINUXBIN)
	chmod +x $(LINUXBIN)

.PHONY: clean
clean:
	rm -f $(EXE) $(EXE)-*-*-*

.PHONY: test
test: $(EXE)
	mkdir -vp $(DIR)
	mv $(EXE) $(DIR)
	go test -v ./...

.PHONY: install
install: $(EXE)
	mkdir -vp ~/bin
	mv $(EXE) ~/bin

.PHONY: docs
docs:
	cd docs; bundle install --path vendor/bundler; bundle exec jekyll build -c _config.yml; cd ..

.PHONY: version
version:
	@echo $(VER)
