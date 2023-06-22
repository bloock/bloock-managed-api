# bloock-managed-api

Managed api is a tool to integrate automatically in your system and include some bloock products:

- Integrity ✅
- Keys
- Authenticity
- Encryption
- Identity
- Availability

See postman doc: https://www.postman.com/bloock/workspace/bloock-api/documentation/16945237-3566feb8-1d4b-45c1-b0a1-f7de0be9348a

## Quickstart

### Requeriments

- Docker Engine(https://docs.docker.com/engine/)
- Makefile toolchain
- Bloock account and API key (https://docs.bloock.com/libraries/authentication/signup)
- Set up a Bloock webhook (https://docs.bloock.com/webhooks/overview)

### Configuration

The service uses viper as a configuration manager, currently supporting environment variables and a YAML configuration file. If no config are provided
we are going to use default values for `BLOOCK_API_HOST`, `BLOOCK_API_PORT`, `BLOOCK_DB_CONNECTION_STRING`, `BLOOCK_ENFORCE_TOLERANCE` and `BLOOCK_API_DEBUG_MODE`.
##### Variables

- BLOOCK_API_PORT: The main API port; default is 8080.
- BLOOCK_API_HOST: The main API host IP; default is 10.0.5.23.
- BLOOCK_API_KEY: Your Bloock API key.
- BLOOCK_WEBHOOK_URL: Your webhook URL.
- BLOOCK_WEBHOOK_SECRET_KEY: Your webhook secret key.
- BLOOCK_ENFORCE_TOLERANCE: Decide if you want to set tolerance when verifying the webhook; true or false.
- BLOOCK_DB_CONNECTION_STRING: Your database URL; e.g., mysql://username:password@localhost:3306/mydatabase.
- BLOOCK_API_DEBUG_MODE: debug mode prints more log information; true or false.

##### Configuration file

The configuration file should be named `config.yaml`. You need to provide the file path as environment variable called BLOOCK_CONFIG_PATH, for example: `BLOOCK_CONFIG_PATH="app/conf/"` (without the filename).
If no config path are provided we try to read from the root directory of the app.
Sample content of `config.yaml`:

```yaml
BLOOCK_DB_CONNECTION_STRING: "file:bloock?mode=memory&cache=shared&_fk=1"
BLOOCK_API_KEY: ""
BLOOCK_WEBHOOK_URL:
BLOOCK_WEBHOOK_SECRET_KEY:  ""
BLOOCK_API_HOST: "0.0.0.0"
BLOOCK_API_PORT: "8080"
BLOOCK_ENFORCE_TOLERANCE: true
BLOOCK_API_DEBUG_MODE: false
```

### Database
The service supports three types of relational databases: MemDB (SQLite), MySQL, and Postgres. You only need to provide the database URL in the following format:

````
MySQL: <user>:<pass>@tcp(<host>:<port>)/<database>?parseTime=True
Postgres: postgresql://username:password@localhost:5432/mydatabase
MemDB: file:dbname?mode=memory&cache=shared&_fk=1
````

### Start the service

#### Locally
The service uses Docker and Makefile to run. There are three example Docker Compose files provided for each database. To start the service locally, follow these steps:

Execute make cache to install the dependencies.
Execute make mocks to create mock files.

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
docker compose -f docker-compose-memdb.yaml up
```
The default Docker Compose file uses MemDB. You can execute the service with `make up` command, which will build the application and start the service in a Docker container.

### Docker image
We push all the service images at: https://hub.docker.com/repository/docker/bloock/managed-api/general

You can use this image in your containerized infrastructure, we recommend to use the highest version.
Also, you can use this image in docker compose file:
```yaml
services:
  bloock-managed-api:
    platform: linux/x86_64
    build:
    image: bloock/managed-api:1.0.0
    container_name: bloock-managed-api
    volumes:
      - ./config.yml:/go/bin/config.yml:ro
    ports:
      - 8080:8080
    environment:
      BLOOCK_DB_CONNECTION_STRING: "file:managed?mode=memory&cache=shared&_fk=1"
      BLOOCK_API_KEY: ""
      BLOOCK_WEBHOOK_URL: ""
      BLOOCK_WEBHOOK_SECRET_KEY: ""
      BLOOCK_API_HOST: "0.0.0.0"
      BLOOCK_API_PORT: "8080"
      BLOOCK_ENFORCE_TOLERANCE: false
      BLOOCK_API_DEBUG_MODE: false

```