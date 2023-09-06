# BLOOCK Managed API

The BLOOCK Managed API is a tool to integrate [BLOOCK](https://bloock.com)'s services using an API instead of our SDKs by taking care of all the integration logic:

---

## Installation

There are two options for running this service:

1. [Docker](#docker-guide)
2. [Standalone](#standalone-guide)

### Docker Guide

Running this API using Docker allows for a quick setup and/or deployment.

#### Locally

To start the service using Docker Compose, you can follow the following steps:

To start the service with MySQL:

```
docker compose -f docker-compose-mysql.yaml up
```

To start the service with Postgres:

```
docker compose -f docker-compose-postgres.yaml up
```

To start the service with MemDB:

```
docker compose -f docker-compose.yaml up
```

> **NOTE:** Remember to update the [configuration](#configuration) variables as required.

### Docker image

We mantain a Docker repository with the latest releases of this repository in [DockerHub](https://hub.docker.com/repository/docker/bloock/managed-api/general).

---

### Standalone Guide

You can also run this service as a common Golang binary if you need it.

#### Standalone Requirements

- Makefile toolchain
- Unix-based operating system (e.g. Debian, Arch, Mac OS X)
- [Go](https://go.dev/) 1.20

#### Standalone Setup

1. Add the required [configuration](#configuration) variables.
2. Run `go run cmd/main.go`

---

### Configuration

The service uses viper as a configuration manager, currently supporting environment variables and a YAML configuration file.

##### Variables

- **BLOOCK_API_HOST**: The API host; default is 0.0.0.0.
- **BLOOCK_API_PORT**: The API port; default is 8080.
- **BLOOCK_API_DEBUG_MODE**: debug mode prints more log information; true or false.
- **BLOOCK_MAX_MEMORY**: The max body size the API will allow; default is 10MB.
- **BLOOCK_DB_CONNECTION_STRING**: Your [database](#database) URL; e.g., mysql://username:password@localhost:3306/mydatabase.
- **BLOOCK_FILE_DIR**: The local directory path where processed files can be stored while waiting for integrity confirmation; default is './'
- **BLOOCK_API_KEY**: Your BLOOCK API key.
- **BLOOCK_WEBHOOK_SECRET_KEY**: Your BLOOCK's webhook secret key.
- **BLOOCK_CLIENT_ENDPOINT_URL**: An endpoint URL where you want to processed file.
- **BLOOCK_AUTHENTICITY_PRIVATE_KEY**: If you want to sign with you own local keys. Here you can set the private key you want to use.
- **BLOOCK_AUTHENTICITY_PUBLIC_KEY**: If you want to sign with you own local keys. Here you can set the public key you want to use.

##### Configuration file

The configuration file should be named `config.yaml`. The service will try to locate this file in the root directory unless the BLOOCK_CONFIG_PATH is defined (i.e. `BLOOCK_CONFIG_PATH="app/conf/"`).

Sample content of `config.yaml`:

```yaml
BLOOCK_API_HOST: "0.0.0.0"
BLOOCK_API_PORT: "8080"
BLOOCK_API_DEBUG_MODE: "false"
BLOOCK_MAX_MEMORY:

BLOOCK_DB_CONNECTION_STRING: "file:bloock?mode=memory&cache=shared&_fk=1"
BLOOCK_FILE_DIR: ""

BLOOCK_API_KEY: ""
BLOOCK_WEBHOOK_SECRET_KEY: ""

BLOOCK_CLIENT_ENDPOINT_URL: ""

BLOOCK_AUTHENTICITY_PRIVATE_KEY: ""
BLOOCK_AUTHENTICITY_PUBLIC_KEY: ""
```

#### Database

The service supports three types of databases: MemDB (SQLite), MySQL, and Postgres. You only need to provide the database URL in the following format:

```
MySQL: <user>:<pass>@tcp(<host>:<port>)/<database>?parseTime=True
Postgres: postgresql://username:password@localhost:5432/mydatabase
MemDB: file:dbname?mode=memory&cache=shared&_fk=1
```

---

## Documentation

You can access the following Postman collection where is the specification for the public endpoint of this API.

[![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)](https://www.postman.com/bloock/workspace/bloock-api/documentation/16945237-3566feb8-1d4b-45c1-b0a1-f7de0be9348a)

---

## License

See [LICENSE](LICENSE.md).
