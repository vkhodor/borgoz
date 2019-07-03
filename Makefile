# Parameters to compile and run application
GOOS?=linux
GOARCH?=amd64

# Current version and commit
VERSION=`git describe --tags`
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags="-X main.Version=$(VERSION)/$(BUILD_TIME)"

all: install

# Compile application
build: borgoz

borgoz:
	@go build $(LDFLAGS)

install: borgoz
	@go install

clean:
	@go clean
