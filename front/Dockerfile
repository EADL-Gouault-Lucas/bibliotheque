FROM golang:1.25.11-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest && templ generate
RUN go build -o server ./cmd/server

FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/web ./web
EXPOSE 3000
CMD ["./server"]
