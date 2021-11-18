.PHONY: build clean test default

VERSION := $(shell grep "const Version " version/version.go | sed -E 's/.*"(.+)"$$/\1/')
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
BINARY_NAME=gl-file-replacer
BINARY_MAC=$(BINARY_NAME)-$(VERSION)-darwin
BINARY_WINDOWS=$(BINARY_NAME)-$(VERSION).exe
ZIP_LINUX=$(BINARY_NAME)-$(VERSION)-amd64-linux.zip
ZIP_WINDOWS=$(BINARY_NAME)-$(VERSION)-amd64-win.zip
ZIP_MAC=$(BINARY_NAME)-$(VERSION)-amd64-darwin.zip

all: test build
default: build
build:
	@echo "building ${BINARY_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags "-X gitlab.com/max-wittig/${BINARY_NAME}/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X gitlab.com/max-wittig/${BINARY_NAME}/version.BuildDate=${BUILD_DATE}" -o bin/${BINARY_NAME}
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X gitlab.com/max-wittig/${BINARY_NAME}/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X gitlab.com/max-wittig/${BINARY_NAME}/version.BuildDate=${BUILD_DATE}" -o bin/${BINARY_NAME}
	cd bin && zip -r $(ZIP_LINUX) $(BINARY_NAME)
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-X gitlab.com/max-wittig/${BINARY_NAME}/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X gitlab.com/max-wittig/${BINARY_NAME}/version.BuildDate=${BUILD_DATE}" -o bin/${BINARY_WINDOWS}
	cd bin && zip -r $(ZIP_WINDOWS) $(BINARY_WINDOWS)
build-mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-X gitlab.com/max-wittig/${BINARY_NAME}/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X gitlab.com/max-wittig/${BINARY_NAME}/version.BuildDate=${BUILD_DATE}" -o bin/${BINARY_MAC}
	cd bin && zip -r $(ZIP_MAC) $(BINARY_MAC)
install: build
	cp bin/${BINARY_NAME} /usr/local/bin/gl-file-replacer
build-all: build-linux build-windows build-mac
clean:
	@test ! -e bin/${BINARY_NAME} || rm -rf bin/
test:
	go test ./...
