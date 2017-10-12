CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-X github.com/Songmu/ghch.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

deps:
	go get ${u} github.com/golang/dep/cmd/dep
	dep ensure

devel-deps: deps
	go get ${u} github.com/golang/lint/golint
	go get ${u} github.com/mattn/goveralls
	go get ${u} github.com/motemen/gobump
	go get ${u} github.com/laher/goxc
	go get ${u} github.com/Songmu/ghch

test: deps
	go test

lint: devel-deps
	go vet
	golint -set_exit_status

cover: devel-deps
	goveralls

build: deps
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/ghch

crossbuild: devel-deps
	goxc -pv=v$(shell gobump show -r) -build-ldflags=$(BUILD_LDFLAGS) \
	  -d=./dist -arch=amd64 -os=linux,darwin,windows \
	  -tasks=clean-destination,xc,archive,rmbin

release:
	_tools/releng

.PHONY: test deps devel-deps lint cover crossbuild release
