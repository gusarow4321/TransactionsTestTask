.PHONY: ent
ent:
	go mod tidy
	go get -d entgo.io/ent/cmd/ent@latest
	go generate ./internal/pkg/ent

.PHONY: build
build:
	mkdir -p bin/ && go build -o ./bin/ ./...