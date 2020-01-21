VERSION = $(shell godzil show-version)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/Songmu/ghch.revision=$(CURRENT_REVISION)"
u := $(if $(update),-u)

export GO111MODULE=on

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: devel-deps
devel-deps:
	sh -c '\
	tmpdir=$$(mktemp -d); \
	cd $$tmpdir; \
	go get ${u} \
	  golang.org/x/lint/golint            \
	  github.com/Songmu/godzil/cmd/godzil \
	  github.com/tcnksm/ghr; \
	rm -rf $$tmpdir'

.PHONY: test
test: deps
	go test

.PHONY: lint
lint: devel-deps
	golint -set_exit_status

.PHONY: build
build: deps
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/ghch

.PHONY: install
install: deps
	go install -ldflags=$(BUILD_LDFLAGS) ./cmd/ghch

.PHONY: release
release:
	godzil release

CREDITS: deps devel-deps go.sum
	godzil credits -w

.PHONY: crossbuild
crossbuild: devel-deps
	godzil crossbuild -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
	  -d=./dist/v$(VERSION) ./cmd/ghch

.PHONY: upload
upload:
	ghr -body="$$(./godzil changelog --latest -F markdown)" v$(VERSION) dist/v$(VERSION)
