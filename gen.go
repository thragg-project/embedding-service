package go_template

// SQL queries generation

// Server code from openapi spec generation
//go:generate oapi-codegen -config api/oapi-codegen/server.yaml api/openapi.yml
//go:generate oapi-codegen -config api/oapi-codegen/spec.yaml api/openapi.yml
//go:generate oapi-codegen -config api/oapi-codegen/types.yaml api/openapi.yml

// Mocks generation
//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=internal/services/embedding.go -destination=internal/services/mocks/embedding_mock.go -package=mocks EmbeddingGateway
//go:generate go run go.uber.org/mock/mockgen@v0.6.0 -source=internal/services/outbox.go -destination=internal/services/mocks/outbox_mock.go -package=mocks EventPublisher
