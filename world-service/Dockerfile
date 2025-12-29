FROM golang:1.25 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

CMD ["./server"]