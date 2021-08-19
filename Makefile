PKG_IMPORT_PATH=github.com/rafaelespinoza/csvtx
BIN_DIR=bin
# If using multiple versions of go, or your target golang version is not in your
# PATH, or whatever reason, specify a path to the golang binary when invoking
# make, for example: make build GO=/path/to/other/golang/bin/go
GO ?= go

.PHONY: all build deps test fmt

all: deps build

build:
	mkdir -pv $(BIN_DIR) && \
		$(GO) build -v -o $(BIN_DIR) \
		-ldflags "\
			-X $(PKG_IMPORT_PATH)/internal/version.BranchName=$(shell git rev-parse --abbrev-ref HEAD) \
			-X $(PKG_IMPORT_PATH)/internal/version.BuildTime=$(shell date --rfc-3339=seconds --utc | tr ' ' 'T') \
			-X $(PKG_IMPORT_PATH)/internal/version.CommitHash=$(shell git rev-parse --short=7 HEAD) \
			-X $(PKG_IMPORT_PATH)/internal/version.GoOSArch=$(shell $(GO) version | awk '{ print $$4 }' | tr '/' '_') \
			-X $(PKG_IMPORT_PATH)/internal/version.GoVersion=$(shell $(GO) version | awk '{ print $$3 }')"

deps:
	$(GO) mod tidy && $(GO) mod download && $(GO) mod verify

# Specify packages to test with P variable. Example:
# make test P='entity repo'
#
# Specify test flags with T variable. Example:
# make test T='-v -count=1 -failfast'
test: P ?= ...
test: pkgpath=$(foreach pkg,$(P),$(shell echo ./internal/$(pkg)))
test:
	$(GO) test $(pkgpath) $(T)

fmt:
	@$(GO) fmt ./... | awk '{ print "Please run go fmt"; exit 1 }'
