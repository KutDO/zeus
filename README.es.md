# Servicio de Descargas

Este es un servicio para descargar archivos con seguimiento de progreso. Está construido utilizando Go, Gorilla Mux para el enrutamiento, y RabbitMQ para la gestión de tareas en cola.

## Requisitos

- Go 1.16 o superior
- RabbitMQ
- Librerías Go adicionales:
  - `github.com/gorilla/mux`
  - `github.com/rabbitmq/amqp091-go`
  - `github.com/BurntSushi/toml`
  - `github.com/swaggo/http-swagger`

## Configuración

1. Asegúrate de tener RabbitMQ corriendo en `localhost:5672`.
2. Crea un archivo `config.toml` en la raíz del proyecto con el siguiente contenido:

    ```toml
    [server]
    address = "0.0.0.0"
    port = 8080
    read_timeout = 15
    write_timeout = 15

    [paths]
    download_directory = "/ruta/a/tu/directorio/de/descargas"
    web_directory = "/ruta/a/tu/directorio/web"

    [rabbitmq]
    uri = "amqp://guest:guest@localhost:5672/"
    queue_name = "download_tasks"
    ```

3. Asegúrate de que los directorios especificados en el archivo `config.toml` existan.

## Instalación

1. Clona el repositorio:

    ```sh
    git clone https://github.com/tu_usuario/tu_repositorio.git
    cd tu_repositorio
    ```

2. Instala las dependencias de Go:

    ```sh
    go get -u github.com/gorilla/mux github.com/rabbitmq/amqp091-go github.com/BurntSushi/toml github.com/swaggo/http-swagger
    ```

3. Compila el proyecto:

    ```sh
    go build -o Zeus
    ```

4. Ejecuta el servidor:

    ```sh
    ./Zeus
    ```

## Uso

### Crear una tarea de descarga

Para crear una nueva tarea de descarga, realiza una solicitud POST a `http://localhost:8080/downloads` con los siguientes parámetros:

- `url`: La URL del archivo a descargar (obligatorio).
- `filename`: Nombre opcional para el archivo descargado.
- `decompress`: Descomprimir el archivo si es `true`.
- `user_id`: ID de usuario opcional para la descarga.

Ejemplo:

```sh
curl -X POST http://localhost:8080/downloads \
  -F "url=https://example.com/file.zip" \
  -F "filename=file.zip" \
  -F "decompress=true" \
  -F "user_id=123e4567-e89b-12d3-a456-426614174000"
