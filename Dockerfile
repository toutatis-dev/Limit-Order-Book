FROM golang:1.25.7 AS builder
WORKDIR /app
COPY . .
RUN go build -o lob ./cmd/

FROM debian:bookworm-slim
COPY --from=builder /app/lob .
CMD ["./lob"]
