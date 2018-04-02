all: container

PREFIX?=simple-test-app
TAG?=1.0
ARCH?=amd64
OS?=linux

server: server.go
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) GOARM=6 go build -a -installsuffix cgo -ldflags '-w -s' -o server

container: server
	docker build $(PREFIX):$(TAG) .

clean:
	rm -f server