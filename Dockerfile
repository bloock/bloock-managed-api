############################
# STEP 1 build executable binary
############################
FROM golang:buster as builder
EXPOSE 8080
RUN update-ca-certificates
WORKDIR /go/bin
COPY . .
RUN CGO_ENABLED=1 go build -buildvcs=false -o /go/bin/managed-api cmd/main.go

############################
# STEP 2 create final image
############################

FROM debian:buster-slim

RUN apt-get update && \
    apt-get install -y ca-certificates sudo && \
    adduser --disabled-password nonroot && \
    echo 'nonroot ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers \
RUN update-ca-certificates

USER nonroot

WORKDIR /home/nonroot/app
COPY --from=builder --chown=nonroot:nonroot /go/bin/managed-api /home/nonroot/app/managed-api
COPY --from=builder --chown=nonroot:nonroot /go/bin/static /home/nonroot/app/static

RUN chmod -R 755 /home/nonroot/app

ENTRYPOINT ["/home/nonroot/app/managed-api"]
