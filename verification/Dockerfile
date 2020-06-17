FROM golang:1.14.3-alpine3.11 as builder
RUN mkdir /app
COPY . /app/
WORKDIR /app
RUN apk update && \
    apk add ca-certificates tzdata
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

#scratch  alpine:latest
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo/
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
