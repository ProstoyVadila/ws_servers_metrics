FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN env GOOS=linux CGO_ENABLED=0 go build -ldflags "-w" -o main

FROM alpine:latest 
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .

CMD ["/app/main"]
