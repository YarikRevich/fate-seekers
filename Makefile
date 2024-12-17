.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: create-local-client
create-local-client: ## Creates fate-seekers-client local directory for API Server
	@mkdir -p $(HOME)/.fate-seekers-client/config
	@mkdir -p $(HOME)/.fate-seekers-client/internal/database

.PHONY: clone-client-config
clone-client-config: create-local-client ## Clones fate-seekers-client configuration files to local directory
	@cp -r ./samples/config/fate-seekers-client/config.yaml $(HOME)/.fate-seekers-client/config

.PHONY: create-local-server
create-local-server: ## Creates fate-seekers-server local directory for API Server
	@mkdir -p $(HOME)/.fate-seekers-server/config

.PHONY: clone-server-config
clone-server-config: create-local-server ## Clones fate-seekers-server configuration files to local directory
	@cp -r ./samples/config/fate-seekers-server/config.yaml $(HOME)/.fate-seekers-server/config

.PHONY: build-server
build-server: clone-server-config ## Builds fate-seekers-server application executable
	@go build -o build/fate-seekers-server ./cmd/fate-seekers-server/...

.PHONY: build-client
build-client: clone-client-config ## Builds fate-seekers-client application executable
	@go build -o build/fate-seekers-client ./cmd/fate-seekers-client/...
