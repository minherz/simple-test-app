all: container

PREFIX?=simple-test-app
TAG?=2.0
ARCH?=amd64
OS?=linux

simple-test-app: server.go
	CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH) GOARM=6 go build -a -installsuffix cgo -ldflags '-w -s' -o simple-test-app

container: simple-test-app
	docker build -t $(PREFIX):$(TAG) .

clean:
	rm -f simple-test-app