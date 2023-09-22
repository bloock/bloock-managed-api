# BLOOCK Managed API

The BLOOCK Managed API is a tool to integrate [BLOOCK](https://bloock.com)'s services using an API instead of our SDKs by taking care of all the integration logic:

---

## Table of Contents

- [Installation](#installation)
    - [Docker Setup Guide](#docker-setup-guide)
        - [Option 1: Pull and Run the Docker Image](#option-1-pull-and-run-the-docker-image)
        - [Option 2: Use Docker Compose with Database Containers](#option-2-use-docker-compose-with-database-containers)
    - [Standalone Setup](#standalone-setup)
        - [Option 3: Clone the GitHub Repository](#option-3-clone-the-github-repository)
- [Configuration](#configuration)
    - [Variables](#variables)
    - [Configuration File](#configuration-file)
- [Database Support](#database-support)
- [Documentation](#documentation)
- [License](#license)

---

## Installation

You have two primary methods to set up and run the Bloock Managed API:

1. [Docker Setup Guide](#docker-setup-guide)
2. [Standalone Setup](#standalone-setup)

Each method has its advantages and use cases.

### Docker Setup Guide

Docker offers a convenient way to package and distribute the API, along with its required dependencies, in a self-contained environment. It's an excellent choice if you want a quick and hassle-free setup, or if you prefer isolation between your application and the host system.

### Option 1: Pull and Run the Docker Image

This option is straightforward and ideal if you want to get started quickly. Follow these steps:

1. **Pull the Docker Image:**

    - Open your terminal or command prompt.

    - Run the following command to pull the Docker image from [DockerHub](https://hub.docker.com/repository/docker/bloock/managed-api/general):

      ```bash
      docker pull bloock/managed-api
      ```

      This command fetches the latest version of the Bloock Managed API image from [DockerHub](https://hub.docker.com/repository/docker/bloock/managed-api/general). We maintain a Docker repository with the latest releases of this repository.


2. **Create a `.env` File:**

    - In your project directory, create a `.env` file. You can use a text editor of your choice to create this file.

    - This file will contain the configuration for the API, including environment variables. Refer to the [Variables](#variables) section for a list of environment variables and their descriptions.

    - In the `.env` file, define the environment variables you want to configure for the API. Each environment variable should be set in the following format:
      ```txt
      VARIABLE_NAME=VALUE
      ```

    - Here's an example of what your `.env` file might look like:

      ```txt
      BLOOCK_DB_CONNECTION_STRING=file:bloock?mode=memory&cache=shared&_fk=1
      BLOOCK_API_KEY=your_api_key
      BLOOCK_WEBHOOK_SECRET_KEY=your_webhook_secret_key
      BLOOCK_CLIENT_ENDPOINT_URL=https://bloock.com/endpoint/to/send/file
      ```

      > **NOTE:** For the **BLOOCK_DB_CONNECTION_STRING** environment variable, you have the flexibility to specify your own MySQL or PostgreSQL infrastructure. Clients can provide their connection string for their database infrastructure. See the [Database](#database-support) section for available connections.

3. **Run the Docker Image with Environment Variables:**

    - Run the following command to start the Bloock Managed API container while passing the `.env` file as an environment variable source:

     ```bash
     docker run --env-file .env -p 8080:8080 bloock/managed-api
     ```

    - This command maps the `.env` file into the container, ensuring that the API reads its configuration from the file. Viper automatically read these environment variables and make them accessible to the code.

    - By default, the above command runs the Docker container in the foreground, displaying API logs and output in your terminal. You can interact with the API while it's running.
    
    3.1. **Running Docker in the Background (Daemon Mode)**

    - Append the `-d` flag to the docker run command as follows:

    ```bash
    docker run -d --env-file config.txt -p 8080:8080 bloock/managed-api
    ```

    The `-d` flag tells Docker to run the container as a background process. You can continue using your terminal for other tasks while the Bloock Managed API container runs silently in the background.


4. **Access the API:**

    - After running the Docker image, the Bloock Managed API will be accessible at `http://localhost:8080`.

    - You can now make API requests to interact with the service.

By following these steps, you can quickly deploy the Bloock Managed API as a Docker container with your customized configuration.

### Option 2: Use Docker Compose with Database Containers

If you need a more complex setup, such as using a specific database like **MySQL**, **Postgres** or **MemDB**, Docker Compose is your choice. Follow these steps:

1. **Choose the Docker Compose File:**

    - In our [repository](https://github.com/bloock/bloock-managed-api), you will find Docker Compose files for different database types:

        - `docker-compose-mysql.yaml` for MySQL
        - `docker-compose-postgres.yaml` for PostgreSQL
        - `docker-compose.yaml for MemDB` (SQLite)


2. **Copy the Chosen Docker Compose File:**

    - Choose the Docker Compose file that corresponds to the database type you want to use. For example, if you prefer MySQL, copy `docker-compose-mysql.yaml`.


3. **Configure Environment Variables:**

    - Open the Docker Compose file in a text editor. Inside the file, locate the environment section for the api service. Here, you can specify environment variables that configure the API.

    - Refer to the [Variables](#variables) section for a list of environment variables and their descriptions.


4. **Set Environment Variables:**

    - In the `environment` section, you can set environment variables in the following format:
      ```yaml
      VARIABLE_NAME: "VALUE"
      ```

    - Here's an example of what your `environment` section might look like:

      ```yaml
      BLOOCK_DB_CONNECTION_STRING: "file:bloock?mode=memory&cache=shared&_fk=1"
      BLOOCK_API_KEY: "your_api_key"
      BLOOCK_WEBHOOK_SECRET_KEY: "your_webhook_secret_key"
      BLOOCK_CLIENT_ENDPOINT_URL: "https://bloock.com/endpoint/to/send/file"
      ```

5. **Run Docker Compose:**

    - Open your terminal, navigate to the directory where you saved the Docker Compose file, and run the following command:

    ```bash
     docker-compose -f docker-compose-mysql.yaml up
     ```

   Replace `docker-compose-mysql.yaml` with the name of the Docker Compose file you selected.

    5.1. **Running Docker in the Background (Daemon Mode)**

    - Append the `-d` flag to the docker run command as follows:

    ```bash
    docker-compose -f docker-compose-mysql.yaml up -d
    ```

    The `-d` flag tells Docker to run the container as a background process. You can continue using your terminal for other tasks while the Bloock Managed API container runs silently in the background.


6. **Access the API:**

    - After running the Docker Compose command, the Bloock Managed API will be accessible at http://localhost:8080. You can make API requests to interact with the service.

By following these steps, you can quickly set up the Bloock Managed API with your chosen database type using the provided Docker Compose files.

### Standalone Setup

Running the API as a standalone application provides more control and flexibility, allowing you to customize and integrate it into your specific environment. Choose this option if you have specific requirements or if you want to modify the API's source code.

### Option 3: Clone the GitHub Repository

You can also run this service as a common Golang binary if you need it.

#### Standalone Requirements

- Makefile toolchain
- Unix-based operating system (e.g. Debian, Arch, Mac OS X)
- [Go](https://go.dev/) 1.20

To deploy the API as a standalone application, follow these steps:

1. **Clone the Repository or Download the Latest Release:**

   1.1. **Clone the Repository:**

    - Open your terminal and navigate to the directory where you want to clone the [repository]((https://github.com/bloock/bloock-identity-managed-api)).

    - Run the following command to clone the [repository]((https://github.com/bloock/bloock-identity-managed-api)):

    ```bash
     git clone https://github.com/bloock/managed-api.git
     ```

   Instead of cloning the repository, it's recommended to download the latest release to ensure you have the most stable and up-to-date version of the Bloock Managed API.

   1.2 **Download the Latest Release:**

    - Visit the [repository's releases page](https://github.com/bloock/bloock-managed-api/releases) on GitHub.

    - Look for the latest release version and select it.

    - Under the Assets section, you will find downloadable files. Choose the appropriate file for your operating system (e.g., Windows, macOS, Linux).

    - Download the selected release file to your local machine.

2. **Navigate to the Repository:**

    - Change your current directory to the cloned repository or downloaded the release file:

    ```bash
     cd managed-api
     ```

3. **Set Up Configuration:**

    - Inside the repository, you'll find a `config.yaml` file.

    - Open `config.yaml` in a text editor and configure the environment variables as needed, following the format described in the [Variables](#variables) section. For example:

    ```yaml
      BLOOCK_DB_CONNECTION_STRING: "file:bloock?mode=memory&cache=shared&_fk=1"
      BLOOCK_API_KEY: "your_api_key"
      BLOOCK_WEBHOOK_SECRET_KEY: "your_webhook_secret_key"
      BLOOCK_CLIENT_ENDPOINT_URL: "https://bloock.com/endpoint/to/send/file"
      ```

4. **Run the Application:**

    - To run the application, execute the following command:

    ```bash
     go run cmd/main.go
     ```

   This command will start the Bloock Managed API as a standalone application, and it will use the configuration provided in the config.yaml file.


5. **Access the API:**

    - After running the application, the Bloock Managed API will be accessible at http://localhost:8080. You can make API requests to interact with the service.

---

## Configuration

The Bloock Managed API leverages Viper, a powerful configuration management library, currently supporting environment variables and a YAML configuration file.

### Variables

Here are the configuration variables used by the Bloock Managed API:

- **BLOOCK_API_KEY** (**REQUIRED**)
    - **Description**: Your unique [BLOOCK API key](https://docs.bloock.com/libraries/authentication/create-an-api-key).
    - **Purpose**: This [API key](https://docs.bloock.com/libraries/authentication/create-an-api-key) is required for authentication and authorization when interacting with the Bloock Identity Managed API. It allows you to securely access and use the API's features.
    - **[Create API Key](https://docs.bloock.com/libraries/authentication/create-an-api-key)**
- **BLOOCK_DB_CONNECTION_STRING** (***OPTIONAL***)
    - **Description**: Your custom database connection URL.
    - **Default**: "file:bloock?mode=memory&cache=shared&_fk=1"
    - **Purpose**: This variable allows you to specify your own [database](#database-support) connection string. You can use it to connect the API to your existing database infrastructure. The format depends on the [database](#database-support) type you choose.
    - **Required**: When docker database container or your existing database infrastructure provided.
- **BLOOCK_WEBHOOK_SECRET_KEY** (***OPTIONAL***)
    - **Description**: Your [BLOOCK webhook secret key](https://docs.bloock.com/webhooks/overview).
    - **Purpose**: The [webhook secret key](https://docs.bloock.com/webhooks/overview) is used to secure and verify incoming webhook requests. It ensures that webhook data is received from a trusted source and has not been tampered with during transmission.
    - **Required**: When you want to certificate data using integrity Bloock product.
    - **[Create webhook](https://docs.bloock.com/webhooks/overview)**
- **BLOOCK_CLIENT_ENDPOINT_URL** (***OPTIONAL***)
    - **Description**: An endpoint URL where you want to send processed files.
    - **Purpose**: This URL specifies the destination where processed files will be sent after successful verification. It can be configured to integrate with other systems or services that require the processed data.
- **BLOOCK_AUTHENTICITY_PRIVATE_KEY** (***OPTIONAL***)
    - **Description**: Private key for signing data.
    - **Purpose**: If you want to sign data using your own local private key, you can specify it here. This private key is used for cryptographic operations to ensure data integrity and authenticity.
- **BLOOCK_AUTHENTICITY_PUBLIC_KEY** (***OPTIONAL***)
    - **Description**: Public key for verifying signed data.
    - **Purpose**: If you're using your own local private key for signing, you should provide the corresponding public key here. The public key is used by others to verify the authenticity of data signed with the private key.
- **BLOOCK_API_HOST** (***OPTIONAL***)
    - **Description**: The API host IP address.
    - **Default**: 0.0.0.0
    - **Purpose**: This variable allows you to specify the IP address on which the Bloock Managed API should listen for incoming requests. You can customize it based on your network configuration.
- **BLOOCK_API_PORT** (***OPTIONAL***)
    - **Description**: The API port number.
    - **Default**: 8080
    - **Purpose**: The API listens on this port for incoming HTTP requests. You can adjust it to match your preferred port configuration.
- **BLOOCK_API_DEBUG_MODE** (***OPTIONAL***)
    - **Description**:  Enable or disable debug mode.
    - **Default**: false
    - **Purpose**: When set to true, debug mode provides more detailed log information, which can be useful for troubleshooting and debugging. Set it to false for normal operation.
- **BLOOCK_FILE_DIR** (***OPTIONAL***)
    - **Description**: The local directory path for storing processed files.
    - **Default**: ./tmp
    - **Purpose**: Processed files can be temporarily stored in this directory while waiting for integrity confirmation. You can configure it to a specific directory path that suits your storage needs.

These configuration variables provide fine-grained control over the behavior of the Bloock Managed API. You can adjust them to match your specific requirements and deployment environment.

### Configuration file

The configuration file should be named `config.yaml`. The service will try to locate this file in the root directory unless the BLOOCK_CONFIG_PATH is defined (i.e. `BLOOCK_CONFIG_PATH="app/conf/"`).

Sample content of `config.yaml`:

```yaml
BLOOCK_API_HOST: "0.0.0.0"
BLOOCK_API_PORT: "8080"
BLOOCK_API_DEBUG_MODE: "false"

BLOOCK_DB_CONNECTION_STRING: "file:bloock?mode=memory&cache=shared&_fk=1"
BLOOCK_FILE_DIR: "./tmp"

BLOOCK_API_KEY: ""
BLOOCK_WEBHOOK_SECRET_KEY: ""
BLOOCK_CLIENT_ENDPOINT_URL: ""

BLOOCK_AUTHENTICITY_PRIVATE_KEY: ""
BLOOCK_AUTHENTICITY_PUBLIC_KEY: ""
```

### Database Support

The Bloock Managed API is designed to be flexible when it comes to database integration. It supports three types of relational databases: **MemDB (SQLite)**, **MySQL**, and **Postgres**. The choice of database type depends on your specific requirements and infrastructure.

Here are the supported database types and how to configure them:

- **MySQL**: To connect to a MySQL database, you can use the following connection string format
   ````
   mysql://user:password@tcp(host:port)/database
   ````

Replace `user`, `password`, `host`, `port`, and `database` with your MySQL database credentials and configuration. This format allows you to specify the MySQL database you want to connect to.

- **Postgres**: For PostgreSQL database integration, use the following connection string format:

   ````
   postgresql://user:password@host/database?sslmode=disable
   ````

Similar to MySQL, replace `user`, `password`, `host`, and `database` with your PostgreSQL database details. Additionally, you can set the `sslmode` as needed. The `sslmode=disable` option is used in the example, but you can adjust it according to your PostgreSQL server's SSL requirements.

- **MemDB (SQLite)**: The API also supports in-memory SQLite databases. To use SQLite, you can specify the connection string as follows:

   ````
   file:dbname?mode=memory&cache=shared&_fk=1
   ````

In this format, `dbname` represents the name of your SQLite database. The API will create an in-memory SQLite database with this name.

If you already have an existing database infrastructure and want to use it with the Bloock Managed API, you have the flexibility to provide your custom database connection string.

`Variable: BLOOCK_DB_CONNECTION_STRING`

The API provides a configuration variable called `BLOOCK_DB_CONNECTION_STRING` that allows you to specify your own database connection string independently of the way you run the API. Whether you run the API as a Docker container or as a standalone application, you can always set this variable to point to your existing database server.

---

## Documentation

You can access the following Postman collection where is the specification for the public endpoint of this API.

[![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)](https://www.postman.com/bloock/workspace/bloock-api/documentation/16945237-1727027f-e3e7-4fe0-969d-afa295eaf2ca)

---

## License

See [LICENSE](LICENSE).

