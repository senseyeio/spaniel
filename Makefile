PACKAGE = github.com/senseyeio/spaniel
GOPACKAGES = $(shell go list ./... | grep -v -e **/*/mock*)

.PHONY: default errcheck fmt lint test tools vet

default: errcheck fmt lint test vet

errcheck:
	@for pkg in $(GOPACKAGES); do errcheck -asserts $$pkg; done

fmt:
	@for pkg in $(GOPACKAGES); do go fmt -x $$pkg; done

lint:
	@for pkg in $(GOPACKAGES); do golint $$pkg; done

test:
	@for pkg in $(GOPACKAGES); do go test -cover $$pkg; done

tools:
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck

vet:
	@for pkg in $(GOPACKAGES); do go vet -x $$pkg; done