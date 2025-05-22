.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/protos/proto/loms/loms.proto

.PHONY: build
build:
	go build -o bin/loms-service cmd/loms/main.go

.PHONY: run
run: build
	./bin/loms-service

.PHONY: test
test:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out