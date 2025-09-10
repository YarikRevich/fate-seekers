.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: generate-proto
generate-proto: ## Generates ProtocolBuffers API for API Server
	@cd api && buf build && buf generate

	@cp api/gen/content/v1/content.pb.go services/fate-seekers-client/pkg/core/networking/content/api && \
	cp api/gen/content/v1/content.pb.go services/fate-seekers-server/pkg/shared/networking/content/api

	@cp api/gen/metadata/v1/metadata_grpc.pb.go services/fate-seekers-client/pkg/core/networking/metadata/api && \
	cp api/gen/metadata/v1/metadata.pb.go services/fate-seekers-client/pkg/core/networking/metadata/api && \
	cp api/gen/metadata/v1/metadata_grpc.pb.go services/fate-seekers-server/pkg/shared/networking/metadata/api && \
	cp api/gen/metadata/v1/metadata.pb.go services/fate-seekers-server/pkg/shared/networking/metadata/api

.PHONY: create-local-client-operational
create-local-client-operational: ## Creates fate-seekers-client operational local directory for API Server
	@mkdir -p $(HOME)/.fate-seekers-client/operational/config
	@mkdir -p $(HOME)/.fate-seekers-client/operational/internal/database

.PHONY: create-local-client-testing
create-local-client-testing: ## Creates fate-seekers-client testing local directory for API Server
	@mkdir -p $(HOME)/.fate-seekers-client/testing/config
	@mkdir -p $(HOME)/.fate-seekers-client/testing/internal/database

.PHONY: clone-client-config-operational
clone-client-config-operational: create-local-client-operational ## Clones fate-seekers-client operational configuration files to local directory
	@cp -r ./samples/config/fate-seekers-client/operational/config.yaml $(HOME)/.fate-seekers-client/operational/config

.PHONY: clone-client-config-testing
clone-client-config-testing: create-local-client-testing ## Clones fate-seekers-client testing configuration files to local directory
	@cp -r ./samples/config/fate-seekers-client/testing/config.yaml $(HOME)/.fate-seekers-client/testing/config

.PHONY: create-local-server
create-local-server: ## Creates fate-seekers-server local directory for API Server
	@mkdir -p $(HOME)/.fate-seekers-server/config
	@mkdir -p $(HOME)/.fate-seekers-server/internal/database

.PHONY: clone-server-config
clone-server-config: create-local-server ## Clones fate-seekers-server configuration files to local directory
	@cp -r ./samples/config/fate-seekers-server/config.yaml $(HOME)/.fate-seekers-server/config

.PHONY: build-server-ui
build-server-ui: clone-server-config ## Builds fate-seekers-server-ui application executable
	@go build -v -tags="server shared" -o build/fate-seekers-server-ui ./cmd/fate-seekers-server-ui/...

.PHONY: build-server-ui-debug
build-server-ui-debug: clone-server-config ## Builds fate-seekers-server-ui application executable in DEBUG mode
	@go build -v -race -tags="server shared" -o build/fate-seekers-server-ui ./cmd/fate-seekers-server-ui/...

.PHONY: build-server-cli
build-server-cli: clone-server-config ## Builds fate-seekers-server-cli application executable
	@go build -v -tags="server shared" -o build/fate-seekers-server-cli ./cmd/fate-seekers-server-cli/...

.PHONY: build-server-cli-debug
build-server-cli-debug: clone-server-config ## Builds fate-seekers-server-cli application executable in DEBUG mode
	@go build -v -race -tags="server shared" -o build/fate-seekers-server-cli ./cmd/fate-seekers-server-cli/...

.PHONY: build-client-operational
build-client-operational: clone-client-config-operational ## Builds fate-seekers-client-operational application executable
	@go build -v -ldflags "-X 'github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config.mode=operational'" -tags="client shared" -o build/fate-seekers-client-operational ./cmd/fate-seekers-client/...

.PHONY: build-client-debug
build-client-operational-debug: clone-client-config ## Builds fate-seekers-client application executable in DEBUG mode
	@go build -v -race -tags="client shared" -o build/fate-seekers-client ./cmd/fate-seekers-client-operational/...

.PHONY: build-client-testing
build-client-testing: clone-client-config-testing ## Builds fate-seekers-client-testing application executable
	@go build -v -ldflags "-X 'github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/config.mode=testing'" -tags="client shared" -o build/fate-seekers-client-testing ./cmd/fate-seekers-client/...