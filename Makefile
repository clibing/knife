BUILD_VERSION   	:= $(shell cat version)
BUILD_DATE      	:= $(shell date "+%F %T")
COMMIT_SHA1     	:= $(shell git rev-parse HEAD)


all: clean
	bash .cross_compile.sh

only: clean
	go build -trimpath -o "dist/knife_darwin_amd64" -ldflags	"-X 'main.version=${BUILD_VERSION}' -X 'main.buildDate=${BUILD_DATE}' -X 'main.commitID=${COMMIT_SHA1}'" 

single: clean
	go build -trimpath -o "dist/knife" -ldflags	"-X 'main.version=${BUILD_VERSION}' -X 'main.buildDate=${BUILD_DATE}' -X 'main.commitID=${COMMIT_SHA1}'" 
	./dist/knife install

release: all
	ghr -u clibing -t ${GITHUB_TOKEN} -replace -recreate -name "Bump ${BUILD_VERSION}" --debug ${BUILD_VERSION} dist

pre-release: all
	ghr -u clibing -t ${GITHUB_TOKEN} -replace -recreate -prerelease -name "Bump ${BUILD_VERSION}" --debug ${BUILD_VERSION} dist

install:
	go install -trimpath -ldflags	"-X 'main.version=${BUILD_VERSION}' \
               						-X 'main.buildDate=${BUILD_DATE}' \
               						-X 'main.commitID=${COMMIT_SHA1}'"
	knife install
debug:
	go install -trimpath -gcflags "all=-N -l"
	knife install

clean:
	rm -rf dist

.PHONY: all release pre-release clean install debug
