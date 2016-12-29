.PHONY: all gofinance test build fmt vet clean release lint errcheck install update release release-check release-copy release-build release-dirs dep_install dep_update

DIST := dist
EXECUTABLE := gofinance

TAGS ?=
TARGETS ?= linux darwin windows
PACKAGES=$(shell go list ./... | grep -v /vendor/)
SOURCES ?= $(shell find . -name "*.go" -type f)
VERSION ?= $(shell git describe --abbrev=0 --tags --always || git rev-parse --short HEAD)
LDFLAGS += -X 'main.Version=$(VERSION)'
LDFLAGS += -X 'main.Githash=$(shell git rev-parse --short HEAD)'
LDFLAGS += -X 'main.Buildstamp=$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')'
LDFLAGS += -X 'main.AppName=$(EXECUTABLE)'

all:build

fmt:
	go fmt $(PACKAGES)

vet:
	go vet $(PACKAGES)

dep_install:
	glide install

dep_update:
	glide up

build:$(EXECUTABLE)
	go build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o bin/${EXECUTABLE} "./cmd/cli"

test:
	for PKG in $(PACKAGES); do go test -v -cover -coverprofile $$GOPATH/src/$$PKG/coverage.txt $$PKG || exit 1; done;

release: release-dirs release-build release-copy release-check

release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

release-build:
	@which gox > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u github.com/mitchellh/gox; \
	fi
	gox -os="$(TARGETS)" -arch="amd64 386" -tags="$(TAGS)" -ldflags="$(EXTLDFLAGS)-s -w $(LDFLAGS)" -output="$(DIST)/binaries/$(EXECUTABLE)-$(VERSION)-{{.OS}}-{{.Arch}}" "./cmd/cli"

release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/release/$(EXECUTABLE)-*),sha1sum $(notdir $(file)) > $(notdir $(file)).sha1;)

errcheck:
	@which errcheck > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u github.com/kisielk/errcheck; \
	fi
	errcheck -verbose $(PACKAGES)

lint:
	@which golint > /dev/null; if [ $$? -ne 0 ]; then \
			go get -u github.com/golang/lint/golint; \
		fi
		for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

clean:
	go clean -x -i ./...
	find . -name coverage.txt -delete
	find . -name *.tar.gz -delete
	find . -name *.db -delete
	-rm -rf bin/* \
		.cover

version:
	@echo $(VERSION)
