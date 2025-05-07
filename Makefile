.PHONY: help
.DEFAULT_GOAL := help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: generate-proto
generate-proto: ## Generates ProtocolBuffers API for API Server
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/metadata/metadata.proto && \
	cp api/metadata/metadata_grpc.pb.go services/fate-seekers-client/pkg/core/networking/metadata/api && \
	cp api/metadata/metadata.pb.go services/fate-seekers-client/pkg/core/networking/metadata/api && \
	cp api/metadata/metadata_grpc.pb.go services/fate-seekers-server/pkg/shared/networking/metadata/api && \
	cp api/metadata/metadata.pb.go services/fate-seekers-server/pkg/shared/networking/metadata/api

	@protoc --go_out=. --go_opt=paths=source_relative api/content/content.proto && \
	cp api/content/content.pb.go services/fate-seekers-client/pkg/core/networking/content/api && \
	cp api/content/content.pb.go services/fate-seekers-server/pkg/shared/networking/content/api

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

.PHONY: build-server-ui
build-server-ui: clone-server-config ## Builds fate-seekers-server-ui application executable
	@go build -v -tags="shared" -o build/fate-seekers-server-ui ./cmd/fate-seekers-server-ui/...

.PHONY: build-server-cli
build-server-cli: clone-server-config ## Builds fate-seekers-server-cli application executable
	@go build -v -o build/fate-seekers-server-cli ./cmd/fate-seekers-server-cli/...

.PHONY: build-client
build-client: clone-client-config ## Builds fate-seekers-client application executable
	@go build -v -tags="client shared" -o build/fate-seekers-client ./cmd/fate-seekers-client/...
