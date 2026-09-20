#Stage 1: Build
FROM golang:1.26-alpine AS builder

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/server ./cmd/main.go

#Stage 2: Runtime
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --from=builder /app/server /server

EXPOSE 8080
CMD ["/server"]