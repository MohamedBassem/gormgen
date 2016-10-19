GO_PKGS=$(foreach pkg, $(shell go list ./...), $(if $(findstring /vendor/, $(pkg)), , $(pkg)))

build: generate build_no_generate

build_no_generate:
	go install $(GO_PKGS)
	go build ./cmd/gormgen

generate:
	go generate

install: generate
	go install $(GO_PKGS)

vet: generate
	go vet $(GO_PKGS)

test: generate
	go test -v

