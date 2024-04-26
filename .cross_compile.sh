#!/usr/bin/env bash

set -e

DIST_PREFIX="knife"
TARGET_DIR="dist"
PLATFORMS="darwin/amd64 darwin/arm64 linux/386 linux/amd64 linux/arm linux/arm64 linux/mips linux/mips64 linux/riscv64 windows/amd64 windows/386 windows/arm windows/arm64"

BUILD_VERSION=$(cat version)
BUILD_DATE=$(date "+%F %T")
COMMIT_SHA1=$(git rev-parse HEAD)

rm -rf ${TARGET_DIR}
mkdir ${TARGET_DIR}

for pl in ${PLATFORMS}; do
    export GOOS=$(echo ${pl} | cut -d'/' -f1)
    export GOARCH=$(echo ${pl} | cut -d'/' -f2)
    export TARGET=${TARGET_DIR}/${DIST_PREFIX}_${GOOS}_${GOARCH}
    if [ "${GOOS}" == "windows" ]; then
        export TARGET=${TARGET_DIR}/${DIST_PREFIX}_${GOOS}_${GOARCH}.exe
    fi
    # export -n CC=x86_64-linux-musl-gcc
    # export -n CXX=x86_64-linux-musl-g++
    # export -n CC=aarch64-linux-musl-gcc
    # export -n CXX=aarch64-linux-musl-g++
    # if [ "${GOOS}" == "linux" ]; then
    #     if [ "${GOARCH}" == "amd64" ]; then
    #         export CC=x86_64-linux-musl-gcc  
    #         export CXX=x86_64-linux-musl-g++
    #         export TARGET=${TARGET_DIR}/${DIST_PREFIX}_${GOOS}_${GOARCH}_musl
    #     fi
    #     if [ "${GOARCH}" == "arm64" ]; then
    #         export CC=aarch64-linux-musl--gcc  
    #         export CXX=aarch64-linux-musl--g++
    #         export TARGET=${TARGET_DIR}/${DIST_PREFIX}_${GOOS}_${GOARCH}_musl
    #     fi
    # fi

    echo "build => ${TARGET}"
    go build -trimpath -o ${TARGET} \
            -ldflags    "-X 'main.version=${BUILD_VERSION}' \
                        -X 'main.buildDate=${BUILD_DATE}' \
                        -X 'main.commitID=${COMMIT_SHA1}'\
                        -w -s"
done

