#!/bin/sh

ARCHS="darwin linux freebsd windows"

if [ $1 == "release" ]; then
    echo "Generating align release binaries..."
    for arch in ${ARCHS}; do
        GOOS=${arch} GOARCH=amd64 go build -v -o bin/align-${arch}
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
        GOOS=freebsd GOARCH=amd64 go build -v -o bin/align-freebsd
        ;;
    "darwin") 
        echo "Building binary for Darwin..."
        GOOS=darwin GOARCH=amd64 go build -v -o bin/align-darwin
        ;;
    "linux") 
        echo "Building binary for Linux..."
        GOOS=linux GOARCH=amd64 go build -v -o bin/align-linux
        ;;
    "windows") 
        echo "Building binary for Windows..."
        GOOS=windows GOARCH=amd64 go build -v -o bin/align-windows.exe
        ;;
esac

exit 0
