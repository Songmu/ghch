VERSION = $(shell godzil show-version)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-X github.com/Songmu/ghch.revision=$(CURRENT_REVISION)"
ifdef update
  u=-u
endif

export GO111MODULE=on

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: devel-deps
devel-deps:
	GO111MODULE=off go get ${u} \
	  golang.org/x/lint/golint            \
	  github.com/mattn/goveralls          \
	  github.com/Songmu/goxz/cmd/goxz     \
	  github.com/Songmu/godzil/cmd/godzil \
	  github.com/Songmu/gocredits/cmd/gocredits

.PHONY: test
test: deps
	go test

.PHONY: lint
lint: devel-deps
	go vet
	golint -set_exit_status

.PHONY: cover
cover: devel-deps
	goveralls

.PHONY: build
build: deps
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/ghch

.PHONY: install
install: deps
	go install -ldflags=$(BUILD_LDFLAGS) ./cmd/ghch

.PHONY: crossbuild
crossbuild: devel-deps
	goxz -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
	  -d=./dist/v$(VERSION) ./cmd/ghch

.PHONY: release
release:
	godzil release

.PHONY: upload
upload:
	ghr v$(VERSION) dist/v$(VERSION)
