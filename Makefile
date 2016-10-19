GO_PKGS=$(foreach pkg, $(shell go list ./...), $(if $(findstring /vendor/, $(pkg)), , $(pkg)))

install:
	go install $(GO_PKGS)

generate:
	go generate

vet:
	go vet $(GO_PKGS)

test: install generate
	go test -v

