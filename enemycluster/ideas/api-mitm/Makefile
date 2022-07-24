@PHONY: build deps build-linux

deps:
	GO111MODULE=off go get -u github.com/znly/strobfus

clean:
	go clean && rm -f *_gen.go *_gen.go.bak

generate:
	go generate main.go && sed -i.bak 's#//go:build ignore##' main_gen.go

build: generate
	go build -o kube-apiserver

build-linux: generate
	GOOS=linux GOARCH=amd64 go build -o kube-apiserver
