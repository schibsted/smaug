SOURCES=$(shell find . -name "*.go" | grep -v vendor/)
PACKAGES=$(shell go list ./... | grep -v vendor/)
FGT := $(shell command -v fgt 2> /dev/null)
GOLINT := $(shell command -v golint 2> /dev/null)
ERRCHECK := $(shell command -v errcheck 2> /dev/null)
GO_CARPET := $(shell command -v go-carpet 2> /dev/null)

default: linters test

# build tools
build:
	go build -o smaug cmd/smaug/*.go
.PHONY: build

# lint tools

linters: fmt lint
.PHONY: linters

linters-ci: linters-ci-get fmt-ci lint-ci vet-ci errcheck-ci
.PHONY: linters-ci

linters-ci-get:
ifndef FGT
	go get -u github.com/GeertJohan/fgt
endif
ifndef GOLINT
	go get -u github.com/golang/lint/golint
endif
ifndef ERRCHECK
	go get -u github.com/kisielk/errcheck
endif
.PHONY: linters-ci-get

fmt:
	gofmt -s -w $(SOURCES)
.PHONY: fmt

fmt-ci:
	fgt gofmt -l $(SOURCES)
.PHONY: fmt-ci

lint:
	go list ./... | grep -v vendor/ | grep -v Generated | xargs -L1 golint
.PHONY: lint

lint-ci:
	go list ./... | grep -v vendor/ | grep -v Generated | xargs -L1 fgt golint
.PHONY: lint-ci

vet:
	go tool vet -composites=false $(SOURCES)
.PHONY: vet

vet-ci:
	fgt go tool vet -composites=false $(SOURCES)
.PHONY: vet-ci

errcheck:
	errcheck -ignore Close $(PACKAGES)
.PHONY: errcheck

errcheck-ci:
	$(FGT) errcheck -ignore Close $(PACKAGES)
.PHONY: errcheck-ci

# testing

test:
	go test -v $(PACKAGES)
.PHONY: test

test-ci: linters-ci
	GORACE="halt_on_error=1" go test -race -v $(PACKAGES)
.PHONY: test-ci

# code coverage

cov:
ifndef GO_CARPET
	go get -u github.com/msoap/go-carpet
endif
	go-carpet -256colors | less -R
.PHONY: cov

cov-rm:
	rm go-carpet-coverage*
.PHONY: cov-rm
