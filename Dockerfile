FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV GOARCH=${GOARCH:-amd64}
ENV CGO_ENABLED=${CGO_ENABLED:-0}
ENV GOOS=${GOOS:-linux}

RUN go build -o server ./cmd/api

# Runtime image
FROM debian:bullseye-slim

COPY --from=builder /app/server .

EXPOSE ${CHALLENGE_SERVER_PORT:-3000}

CMD ["/server"]
