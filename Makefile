default: help

.PHONY: help build install lint fmt test testacc gendocs genresources genall
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build the provider"
	@echo "  install      Install the provider"
	@echo "  lint         Run linter"
	@echo "  fmt          Run gofmt"
	@echo "  test         Run tests"
	@echo "  testacc      Run acceptance tests"
	@echo "  gendocs      Generate docs"
	@echo "  genresources Generate resources"
	@echo "  genall       Generate all"

build:
	go build -v ./...

install: build
	go install -v ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w -e .

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

gendocs:
	terraform fmt -recursive examples/
	go tool tfplugindocs generate --provider-name azname

genresources:
	cd internal/resources; go generate ./...

genall: gendocs genresources