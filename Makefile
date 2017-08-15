ifdef update
  u=-u
endif

test: test-deps
	go test

deps:
	go get ${u} -d -v ./...

test-deps:
	go get ${u} -d -t -v ./...

devel-deps: test-deps
	go get ${u} github.com/golang/lint/golint
	go get ${u} github.com/mattn/goveralls
	go get ${u} github.com/motemen/gobump
	go get ${u} github.com/laher/goxc

lint: test-deps
	go vet ./...
	golint -set_exit_status ./...

cover: devel-deps
	goveralls

crossbuild: devel-deps
	goxc -pv=v$(shell gobump show -r) -d=./dist -arch=amd64 -os=linux,darwin,windows -tasks=clean-destination,xc,archive,rmbin

release:
	_tools/releng

.PHONY: test deps test-deps devel-deps lint cover crossbuild release
