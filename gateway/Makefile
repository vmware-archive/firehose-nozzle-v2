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
	rm -f Gopkg.lock

HAS_DEP := $(shell command -v dep;)

boostrap: bootstrap

.PHONY: bootstrap
bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/maxbrunsfeld/counterfeiter
endif
	dep ensure -v

all: fmt test
