
### README en Inglés

```markdown
# Download Service

This is a service for downloading files with progress tracking. It is built using Go, Gorilla Mux for routing, and RabbitMQ for task queue management.

## Requirements

- Go 1.16 or higher
- RabbitMQ
- Additional Go libraries:
  - `github.com/gorilla/mux`
  - `github.com/rabbitmq/amqp091-go`
  - `github.com/BurntSushi/toml`
  - `github.com/swaggo/http-swagger`

## Configuration

1. Ensure you have RabbitMQ running on `localhost:5672`.
2. Create a `config.toml` file in the project root with the following content:

    ```toml
    [server]
    address = "0.0.0.0"
    port = 8080
    read_timeout = 15
    write_timeout = 15

    [paths]
    download_directory = "/path/to/your/download/directory"
    web_directory = "/path/to/your/web/directory"

    [rabbitmq]
    uri = "amqp://guest:guest@localhost:5672/"
    queue_name = "download_tasks"
    ```

3. Ensure the directories specified in the `config.toml` file exist.

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/your_username/your_repository.git
    cd your_repository
    ```

2. Install Go dependencies:

    ```sh
    go get -u github.com/gorilla/mux github.com/rabbitmq/amqp091-go github.com/BurntSushi/toml github.com/swaggo/http-swagger
    ```

3. Build the project:

    ```sh
    go build -o Zeus
    ```

4. Run the server:

    ```sh
    ./Zeus
    ```

## Usage

### Create a Download Task

To create a new download task, make a POST request to `http://localhost:8080/downloads` with the following parameters:

- `url`: The URL of the file to download (required).
- `filename`: Optional name for the downloaded file.
- `decompress`: Decompress the file if `true`.
- `user_id`: Optional user ID for the download.

Example:

```sh
curl -X POST http://localhost:8080/downloads \
  -F "url=https://example.com/file.zip" \
  -F "filename=file.zip" \
  -F "decompress=true" \
  -F "user_id=123e4567-e89b-12d3-a456-426614174000"
