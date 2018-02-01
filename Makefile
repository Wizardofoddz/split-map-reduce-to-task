# Go parameters
GOCMD   := go
DEPCMD  := dep
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST  := $(GOCMD) test
GOGET   := $(GOCMD) get
GORELEASE := env GCO_ENABLED=0 $(GOBUILD)
BINNAME := split-map-reduce-to-task

all: deps build
build:
	$(GOBUILD)
release:
	# builds current os release builds with license enforcement
	$(GORELEASE) -o "./dist/bin/$(BINNAME)" -a -ldflags

create-release-dir:
	mkdir -p ./dist/release

crossrelease: release-darwin release-linux release-windows
	# builds releases for all target operating systems

release-darwin: release-darwin-amd64

release-darwin-amd64: create-release-dir
	env GOOS=darwin GOARCH=amd64 $(GORELEASE) -o "./dist/bin/darwin-amd64/$(BINNAME)"
	tar -cz -C ./dist/bin/darwin-amd64 -f ./dist/release/$(BINNAME)-darwin-amd64.tar.gz $(BINNAME)

release-linux: release-linux-amd64 release-linux-arm release-linux-386

release-linux-amd64: create-release-dir
	env GOOS=linux GOARCH=amd64 $(GORELEASE) -o "./dist/bin/linux-amd64/$(BINNAME)"
	tar -cz -C ./dist/bin/linux-amd64 -f ./dist/release/$(BINNAME)-linux-amd64.tar.gz $(BINNAME)

release-linux-arm: create-release-dir
	env GOOS=linux GOARCH=arm $(GORELEASE) -o "./dist/bin/linux-arm/$(BINNAME)"
	tar -cz -C ./dist/bin/linux-arm -f ./dist/release/$(BINNAME)-linux-arm.tar.gz $(BINNAME)

release-linux-386: create-release-dir
	env GOOS=linux GOARCH=386 $(GORELEASE) -o "./dist/bin/linux-386/$(BINNAME)"
	tar -cz -C ./dist/bin/linux-386 -f ./dist/release/$(BINNAME)-linux-386.tar.gz $(BINNAME)

release-windows: release-windows-amd64 release-windows-386

release-windows-amd64: create-release-dir
	env GOOS=windows GOARCH=amd64 $(GORELEASE) -o "./dist/bin/windows-amd64/$(BINNAME)"
	tar -cz -C ./dist/bin/windows-amd64 -f ./dist/release/$(BINNAME)-windows-amd64.tar.gz $(BINNAME)

release-windows-386: create-release-dir
	env GOOS=windows GOARCH=386 $(GORELEASE) -o "./dist/bin/windows-386/$(BINNAME)"
	tar -cz -C ./dist/bin/windows-386 -f ./dist/release/$(BINNAME)-windows-386.tar.gz $(BINNAME)

clean:
	$(GOCLEAN)
	rm -rf dist

deps:
	$(DEPCMD) ensure
