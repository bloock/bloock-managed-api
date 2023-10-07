############################
# STEP 1 build executable binary
############################
FROM golang:buster as builder
EXPOSE 8080
RUN update-ca-certificates
WORKDIR /go/bin
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -a -o /go/bin/managed-api cmd/main.go

############################
# STEP 2 create final image
############################

FROM debian:buster-slim

RUN apt-get update && apt-get install -y ca-certificates
RUN update-ca-certificates

WORKDIR /app
COPY --from=builder /go/bin/managed-api /app/managed-api

ENTRYPOINT ["/app/managed-api"]
