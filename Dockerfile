ARG  BUILDER_IMAGE=golang:1.20
FROM ${BUILDER_IMAGE}
EXPOSE 8080
WORKDIR /go/bin
COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -a -o /go/bin/run cmd/main.go

CMD ["run"]

