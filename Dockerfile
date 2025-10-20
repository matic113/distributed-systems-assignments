FROM golang:1.25.3-alpine3.22 AS builder
WORKDIR /app

COPY src/go.mod  ./
RUN go mod download

COPY src/ ./
RUN go build -o /app/main .

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .

# making a volume for saved_users folder to persist data
VOLUME /app/saved_users

EXPOSE 8080
CMD ["./main"]