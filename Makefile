# Parameters to compile and run application
GOOS?=linux
GOARCH?=amd64

# Current version and commit
VERSION=`git describe --tags`
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags="-X main.Version=$(VERSION)/$(BUILD_TIME)"

all: build

# Compile application
build: 
	@go build $(LDFLAGS)
