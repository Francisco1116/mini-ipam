FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o mini-ipam main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/mini-ipam .
EXPOSE 8080
CMD ["./mini-ipam"]