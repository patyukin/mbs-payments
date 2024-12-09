FROM golang:1.23.2-alpine3.20 AS builder

COPY . /app
WORKDIR /app

RUN go mod download
RUN go mod tidy
RUN go build -o ./bin/payment cmd/payment/main.go

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /app/bin/payment .
ENV YAML_CONFIG_FILE_PATH=config.yaml
COPY config.yaml config.yaml

CMD ["./payment"]
