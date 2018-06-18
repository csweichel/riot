# ########################################################## #
# Makefile for Golang Project
# Includes cross-compiling, installation, cleanup
#
# Origin: https://gist.github.com/cjbarker/5ce66fcca74a1928a155cfb3fea8fac4
# ########################################################## #

# Check for required command tools to build or stop immediately
EXECUTABLES = git go find pwd
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=riot
VERSION=0.1.0
BUILD=`git rev-parse HEAD`
PLATFORMS=darwin linux
ARCHITECTURES=386 amd64 arm
GOARM=6

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X cmd.Version=${VERSION} -X cmd.Build=${BUILD}"

default: build

all: clean build_all install

build:
	go get -v ./...
	go build ${LDFLAGS} -o ${BINARY}

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); export GOARM=$(GOARM); go build -v -o $(BINARY)-$(GOOS)-$(GOARCH))))

install:
	go install ${LDFLAGS}

# Remove only what we've created
clean:
	find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

.PHONY: check clean build_all all