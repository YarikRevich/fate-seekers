.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build-server
build-server: ## Builds fate-seekers-server application executable
	@go build -o build/fate-seekers-server ./cmd/fate-seekers-server/...

.PHONY: build-client
build-client: ## Builds fate-seekers-client application executable
	@go build -o build/fate-seekers-client ./cmd/fate-seekers-client/...
