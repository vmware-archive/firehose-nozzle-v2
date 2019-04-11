default: all

GO_PACKAGES = $$(go list ./... | grep -v vendor)
GO_FILES = $$(find . -name "*.go" | grep -v vendor | uniq)

unit-test:
	@go test ${GO_PACKAGES}

fmt:
	gofmt -s -l -w $(GO_FILES)

vet:
	@go vet ${GO_PACKAGES}

test: generate unit-test vet

generate:
	go generate ./...

cleandep:
	rm -rf vendor
	rm -f go.sum

HAS_CFER := $(shell command -v counterfeiter;)

boostrap: bootstrap
.PHONY: bootstrap
bootstrap:
ifndef HAS_CFER
	go get -u github.com/maxbrunsfeld/counterfeiter
endif
	go mod vendor

all: fmt test
