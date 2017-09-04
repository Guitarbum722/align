OS := $(shell uname)

build: clean
	cd cmd/align
	go build -o bin/align

test:
	go test -v

clean:
	go clean
	rm -f align
	rm -f bin/*

install: clean
ifeq ($(OS),Darwin)
	./build.sh darwin
	cp -f bin/align-darwin /usr/local/bin/align
endif 
ifeq ($(OS),Linux)
	./build.sh linux
	cp -f bin/align-linux /usr/local/bin/align
endif
ifeq ($(OS),FreeBSD)
	./build.sh freebsd
	cp -f bin/align-freebsd /usr/local/bin/align
endif
uninstall: 
	rm -f /usr/local/bin/align

release:
	./build.sh release
