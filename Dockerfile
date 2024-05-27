############################
# STEP 1 build executable binary
############################
FROM golang:1.22.3-bookworm AS build

RUN useradd -u 1001 nonroot
RUN update-ca-certificates

WORKDIR /go/bin
COPY . .
RUN CGO_ENABLED=1 go build -buildvcs=false -ldflags="-linkmode external -extldflags -static" -o /go/bin/managed-api cmd/main.go


############################
# STEP 2 create final image
############################
FROM debian:bookworm-slim

WORKDIR /go/bin

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/bin/managed-api managed-api
COPY --from=build /go/bin/static static

USER nonroot
EXPOSE 8080

CMD ["/go/bin/managed-api"]
