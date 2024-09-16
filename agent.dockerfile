# Этап, на котором выполняется сборка приложения
FROM golang:1.22.0-alpine as builder
RUN apk update && apk add --no-cache git
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /main cmd/agent/main.go

# Финальный этап, копируем собранное приложение
FROM alpine:3
RUN apk add curl
COPY --from=builder main /bin/main
ENTRYPOINT ["/bin/main"]