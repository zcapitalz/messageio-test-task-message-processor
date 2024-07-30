.PHONY: generate-swagger-docs
generate-server-swagger-docs:
	swag init --generalInfo=./internal/server/server.go --parseInternal --parseDependency --output=./api/server

.PHONY: build-processor-image
build-processor-image:
	docker build \
		--file ./Dockerfile-processor \
		--tag message-processor:latest \
		.