test: deps
	go test

deps:
	go get -d -v -t ./...
	go get github.com/golang/lint/golint
	go get github.com/mattn/goveralls

LINT_RET = .golint.txt
lint: deps
	go vet ./...
	rm -f $(LINT_RET)
	golint ./... | tee $(LINT_RET)
	test ! -s $(LINT_RET)

cover: deps
	goveralls

release:
	_tools/releng

.PHONY: test deps lint cover release
