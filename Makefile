CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-X github.com/Songmu/ghch.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

GO ?= GO111MODULE=on go

devel-deps:
	GO111MODULE=off go get ${u} \
	  golang.org/x/lint/golint             \
	  github.com/mattn/goveralls           \
	  github.com/motemen/gobump/cmd/gobump \
	  github.com/Songmu/goxz/cmd/goxz      \
	  github.com/Songmu/ghch/cmd/ghch

test:
	$(GO) test

lint: devel-deps
	$(GO) vet
	golint -set_exit_status

cover: devel-deps
	goveralls

build:
	$(GO) build -ldflags=$(BUILD_LDFLAGS) ./cmd/ghch

crossbuild: devel-deps
	GO111MODULE=on goxz -pv=v$(shell gobump show -r) -build-ldflags=$(BUILD_LDFLAGS) \
	  -d=./dist/v$(shell gobump show -r) ./cmd/ghch

release:
	_tools/releng

.PHONY: test devel-deps lint cover build crossbuild release
