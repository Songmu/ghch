VERSION = $(shell godzil show-version)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-X github.com/Songmu/ghch.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

export GO111MODULE=on

deps:
	go get ${u} -d

devel-deps:
	GO111MODULE=off go get ${u} \
	  golang.org/x/lint/golint             \
	  github.com/mattn/goveralls           \
	  github.com/Songmu/goxz/cmd/goxz      \
	  github.com/Songmu/godzil/cmd/godzil

test: deps
	go test

lint: devel-deps
	go vet
	golint -set_exit_status

cover: devel-deps
	goveralls

build:
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/ghch

crossbuild: devel-deps
	goxz -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
	  -d=./dist/v$(VERSION) ./cmd/ghch

release:
	godzil release

upload:
	ghr v$(VERSION) dist/v$(VERSION)

.PHONY: deps devel-deps test lint cover build crossbuild release upload
