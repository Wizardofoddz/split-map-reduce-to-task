# Go parameters
GOCMD   := go
DEPCMD  := dep
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST  := $(GOCMD) test
GOGET   := $(GOCMD) get

DISTDIR := $(PWD)/dist

BINOS := -$(GOOS)
BINARCH := -$(GOARCH)
BINTAG := -$(CIRCLE_TAG)

ifeq ($(GOOS), windows)
	BINEXT := .exe
endif
ifeq ($(GOOS),)
	BINOS :=
endif
ifeq ($(GOARCH),)
	BINARCH :=
endif
ifeq ($(CIRCLE_TAG),)
	BINTAG :=
endif

# Binary names
BINNAME    := split-map-reduce-to-task$(BINTAG)$(BINOS)$(BINARCH)
BINNAMEEXT := $(BINNAME)$(BINEXT)

# Release params
VERSION         := $(shell tr -d '\n' < ./VERSION)
RELEASE_LDFLAGS := -s -X main.Version=$(VERSION)
GORELEASE       := env GOOS=$(GOOS) GOARCH=$(GOARCH) GCO_ENABLED=0 $(GOBUILD) -a -ldflags "$(RELEASE_LDFLAGS)"

all: deps build
build:
	$(GOBUILD)

release:
	$(GORELEASE) -o "$(DISTDIR)/bin/$(GOOS)-$(GOARCH)/$(BINNAMEEXT)"
	tar -cz -C $(DISTDIR)/bin/$(GOOS)-$(GOARCH) -f $(DISTDIR)/release/$(BINNAME).tar.gz $(BINNAMEEXT)
	zip -Dj9 $(DISTDIR)/release/$(BINNAME).zip $(DISTDIR)/bin/$(GOOS)-$(GOARCH)/$(BINNAMEEXT)

create-release-dir:
	mkdir -p $(DISTDIR)/release

crossrelease: release-darwin release-linux release-windows
	# builds releases for all target operating systems

release-darwin: release-darwin-amd64

release-darwin-amd64: create-release-dir
	env GOOS=darwin GOARCH=amd64 $(MAKE) release

release-linux: release-linux-amd64 release-linux-arm release-linux-arm64 release-linux-386

release-linux-amd64: create-release-dir
	env GOOS=linux GOARCH=amd64 $(MAKE) release

release-linux-arm: create-release-dir
	env GOOS=linux GOARCH=arm $(MAKE) release

release-linux-arm64: create-release-dir
	env GOOS=linux GOARCH=arm64 $(MAKE) release

release-linux-386: create-release-dir
	env GOOS=linux GOARCH=386 $(MAKE) release

release-windows: release-windows-amd64 release-windows-386

release-windows-amd64: create-release-dir
	env GOOS=windows GOARCH=amd64 $(MAKE) release

release-windows-386: create-release-dir
	env GOOS=windows GOARCH=386 $(MAKE) release

clean:
	$(GOCLEAN)
	rm -rf dist

deps:
	$(DEPCMD) ensure
