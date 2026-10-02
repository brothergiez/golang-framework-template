# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG SERVICE_TYPE=api
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/service ./cmd/${SERVICE_TYPE}

FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser
WORKDIR /app
COPY --from=builder /out/service /app/service
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/service"]
