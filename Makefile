default: help

# Run acceptance tests
.PHONY: help testacc docs
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  testacc  - Run acceptance tests"
	@echo "  docs  - Generate documentation for the provider"

testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

docs:
	terraform fmt -recursive ./examples/
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate -provider-name azname