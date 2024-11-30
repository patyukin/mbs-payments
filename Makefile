.PHONY: start stop rebuild gen up down restart deps tidy test

LOCAL_BIN:=$(CURDIR)/bin

up:
	docker compose up -d

down:
	docker compose down

start:
	docker compose start

stop:
	docker compose stop

restart:
	docker compose restart

rebuild:
	docker compose down -v --remove-orphans
	docker compose up -d --build

gen:
	make install-deps
	make gen-api

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v1.1.0
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

tidy:
	go mod tidy

test:
	go test ./...

gen-api:
	mkdir -p pkg/proto/auth_v1
	protoc --proto_path api/auth_v1 \
	--go_out=pkg/proto/auth_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/proto/auth_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/proto/auth_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	api/auth_v1/auth.proto
	mkdir -p pkg/proto/payment_v1
	protoc --proto_path api/payment_v1 \
	--go_out=pkg/proto/payment_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/proto/payment_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	--validate_out lang=go:pkg/proto/payment_v1 --validate_opt=paths=source_relative \
	--plugin=protoc-gen-validate=bin/protoc-gen-validate \
	api/payment_v1/payment.proto
