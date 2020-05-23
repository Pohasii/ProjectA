FROM golang:1.14.3-alpine3.11
RUN mkdir /app
COPY . /app/
WORKDIR /app
RUN apk update && \
    apk add ca-certificates tzdata
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o server .

#scratch  alpine:latest
FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /usr/share/zoneinfo /usr/share/zoneinfo/
WORKDIR /root/
COPY --from=0 /app/server .
CMD ["./server"]
