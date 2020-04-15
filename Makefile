# MAKEFILE
#
# @author      Ollie Parsley <ollie@ollieparsley.com>
# @link        https://github.com/ollieparsley/twitter-api-metrics
# ------------------------------------------------------------------------------

# Ensure everyone is using bash. Note that Ubuntu now uses dash which doesn't support PIPESTATUS.
SHELL=/bin/bash

# Project version
MINOR_VERSION = 0
ifdef DRONE_BUILD_NUMBER
  MINOR_VERSION = $(DRONE_BUILD_NUMBER)
endif
VERSION=$(shell cat VERSION).$(MINOR_VERSION)

# Name of package
PKGNAME=twitter-api-metrics

# Description
DESCRIPTION=Fetch your Twitter API rate limits

# Destdir
DESTDIR?=target/build/root/

# Binary path (where the executable files will be installed)
BINPATH=usr/bin/

# Docker
DOCKER_USERNAME?=ollieparsley
DOCKER_PASSWORD?=not_this_password

# Environment variable exports
export GO111MODULE=on

# Chart packaging
version:
	echo $(VERSION) > VERSION

# Init
init:
	go mod init

# Get the dependencies
deps:
	go get ./...
	go get golang.org/x/lint/golint

# Run the unit tests
test:
	go test -v ./...

# Check for syntax errors
vet:
	go vet ./...

# Go fmt
fmt:
	go fmt ./...

# Check for style errors
lint:
	golint -set_exit_status ./...

# Alias to run targets: fmtcheck test vet lint
qa: fmt test vet lint

# Compile the application
build: deps
	@mkdir -p $(DESTDIR)
	@mkdir -p $(DESTDIR)$(BINPATH)
	go build -o $(DESTDIR)$(BINPATH)$(PKGNAME) ./main.go
