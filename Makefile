PKG_IMPORT_PATH=github.com/rafaelespinoza/csvtx
BIN_DIR=bin
# If using multiple versions of go, or your target golang version is not in your
# PATH, or whatever reason, specify a path to the golang binary when invoking
# make, for example: make build GO=/path/to/other/golang/bin/go
GO ?= go
GOSEC ?= gosec
SRC_PATHS=$(PKG_IMPORT_PATH) $(PKG_IMPORT_PATH)/internal/...

.PHONY: all build deps test fmt vet

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
	$(GO) mod tidy

# Specify packages to test with PKGS variable. Example:
# make test PKGS='./internal/entity ./internal/cmd'
#
# Specify test flags with FLAGS variable. Example:
# make test FLAGS='-v -count=1 -failfast'
test: PKGS ?= ./...
test: pkgpaths=$(shell $(GO) list $(PKGS))
test:
	$(GO) test $(pkgpaths) $(FLAGS)

fmt:
	@$(GO) fmt $(SRC_PATHS) | awk '{ print "Please run go fmt"; exit 1 }'

vet:
	$(GO) vet $(FLAGS) $(SRC_PATHS)

# Run a security scanner over the source code. This Makefile won't install the
# scanner binary for you, so check out the gosec README for instructions:
# https://github.com/securego/gosec
#
# If necessary, specify the path to the built binary with the GOSEC env var.
#
# Also note, the package paths (last positional input to gosec command) should
# be a "relative" package path. That is, starting with a dot.
gosec:
	$(GOSEC) $(FLAGS) . ./internal/...
