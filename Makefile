PACKAGE=github.com/genshen/wssocks

.PHONY: clean all


LDFLAGS= -v -ldflags "-X 'github.com/genshen/wssocks/version.buildHash=`git rev-parse HEAD`' \
 -X 'github.com/genshen/wssocks/version.buildTime=`date`' \
 -X 'github.com/genshen/wssocks/version.buildGoVersion=`go version | cut -f 3,4 -d" "`'"

all: wssocks-sangfor-linux-amd64 wssocks-sangfor-linux-arm64 wssocks-sangfor-darwin-amd64 wssocks-sangfor-darwin-arm64 wssocks-sangfor-windows-amd64.exe

wssocks-sangfor-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o wssocks-sangfor-linux-amd64 ${PACKAGE}

wssocks-sangfor-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o wssocks-sangfor-linux-arm64 ${PACKAGE}

wssocks-sangfor-darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o wssocks-sangfor-darwin-arm64 ${PACKAGE}

wssocks-sangfor-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o wssocks-sangfor-darwin-amd64 ${PACKAGE}

wssocks-sangfor-windows-amd64.exe:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o wssocks-sangfor-windows-amd64.exe ${PACKAGE}

wssocks :
	go build -o wssocks

clean:
	rm -f wssocks-sangfor-linux-amd64 wssocks-sangfor-linux-arm64 wssocks-sangfor-darwin-arm64 wssocks-sangfor-darwin-amd64 wssocks-sangfor-windows-amd64.exe
