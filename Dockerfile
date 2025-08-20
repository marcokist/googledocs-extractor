FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/doc-extractor-api

FROM gcr.io/distroless/static-debian11 AS final

WORKDIR /app

COPY --from=builder /app/server .

COPY credentials.json .

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/server"]