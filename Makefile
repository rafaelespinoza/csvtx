# If using multiple versions of go, or your target golang version is not in your
# PATH, or whatever reason, specify a path to the golang binary when invoking
# make, for example: make build GO=/path/to/other/golang/bin/go
GO ?= go

.PHONY: all build deps test fmt

all: deps build

build:
	$(GO) build ./...

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
