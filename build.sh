#!/bin/sh

ARCHS="darwin linux freebsd windows"
CMD_DIR="cmd/align"
BUILD_CMD="go build -v -o"

if [ $1 == "release" ]; then
    echo "Generating align release binaries..."
    for arch in ${ARCHS}; do
        GOOS=${arch} GOARCH=amd64 ${BUILD_CMD} bin/align-${arch}
    done
fi

case "$1" in
    "release") 
        echo "Building release..."
        for arch in ${ARCHS}; do
            GOOS=${arch} GOARCH=amd64 go build -v -o bin/align-${arch}
            tar -czvf bin/align-${arch}.tar.gz bin/align-${arch}
        done
        ;;
    "freebsd") 
        echo "Building binary for FreeBSD..."
        cd ${CMD_DIR}
        GOOS=freebsd GOARCH=amd64 ${BUILD_CMD} ../../bin/align-freebsd
        ;;
    "darwin") 
        echo "Building binary for Darwin..."
        cd ${CMD_DIR}
        GOOS=darwin GOARCH=amd64 ${BUILD_CMD} ../../bin/align-darwin
        ;;
    "linux") 
        echo "Building binary for Linux..."
        cd ${CMD_DIR}
        GOOS=linux GOARCH=amd64 ${BUILD_CMD} ../../bin/align-linux
        ;;
    "windows") 
        echo "Building binary for Windows..."
        cd ${CMD_DIR}
        GOOS=windows GOARCH=amd64 ${BUILD_CMD} ../../bin/align-windows.exe
        ;;
esac

exit 0
