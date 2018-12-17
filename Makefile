
GIT_SHA=$(shell git rev-parse HEAD)
DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_INFO_IMPORT_PATH=night-watch/pkg/build
BUILD_INFO=-ldflags "-X $(BUILD_INFO_IMPORT_PATH).commitSHA=$(GIT_SHA) -X $(BUILD_INFO_IMPORT_PATH).date=$(DATE)"



.PHONY: build
build_linux: build
	CGO_ENABLED=0 GOOS=linux installsuffix=cgo go build $(BUILD_INFO) -o ./build/nightwatch ./main.go

build_mac: build
	GOOS=darwin go build $(BUILD_INFO) -o ./build/nightwatch-mac ./main.go

