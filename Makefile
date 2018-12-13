
GIT_SHA=$(shell git rev-parse HEAD)
DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_INFO_IMPORT_PATH=github.com/dashbase/watch-append/pkg/build
BUILD_INFO=-ldflags "-X $(BUILD_INFO_IMPORT_PATH).commitSHA=$(GIT_SHA) -X $(BUILD_INFO_IMPORT_PATH).date=$(DATE)"



.PHONY: build
build_linux: build
	CGO_ENABLED=0 GOOS=linux installsuffix=cgo go build -o ./bin/nightwatch ./main.go
