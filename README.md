# Go URL Shortener

[![Go Report Card](https://goreportcard.com/badge/github.com/H0llyW00dzZ/go-urlshortner)](https://goreportcard.com/report/github.com/H0llyW00dzZ/go-urlshortner)
[![Go Reference](https://pkg.go.dev/badge/github.com/H0llyW00dzZ/go-urlshortner.svg)](https://pkg.go.dev/github.com/H0llyW00dzZ/go-urlshortner)
[![CI: CodeQL Unit Testing Advanced](https://github.com/H0llyW00dzZ/go-urlshortner/actions/workflows/CodeQL.yml/badge.svg)](https://github.com/H0llyW00dzZ/go-urlshortner/actions/workflows/CodeQL.yml)

This is a simple URL shortening service written in Go, which utilizes the Gin Web Framework for routing and Google Cloud Datastore for persistent storage.

## Introduction

<p align="center">
  <img src="https://i.imgur.com/T04JNPd.jpg" alt="Go Picture">
</p>

This project aims to provide a straightforward and scalable approach to creating short aliases for lengthy URLs. It's constructed in Go and is designed to be simple to deploy and maintain. The service includes basic functionalities such as generating a shortened URL and redirecting to the original URL when accessed.

## Features

- **🔗 Shorten URLs**: The service allows users to convert long URLs into short, manageable identifiers that redirect to the original URLs. This makes the links easier to share on platforms with character limits or for readability purposes.

- **↪️ Redirection**: When a client accesses a short link, the service performs an HTTP or HTTPS redirect to the original, long URL associated with the short identifier.

- **🔌 Simple Integration**: The service exposes RESTful API endpoints for creating, retrieving, updating, and deleting shortened URLs, which can be easily used by external applications or services.

- **🏴‍☠️ Dependency Injection (Greetz Thx to: Gopher Unit Hacker)**: The application's architecture includes dependency injection for core components like the logger and the Datastore client. This design pattern enhances testability and maintainability by decoupling the components and their dependencies.

- **⚖️ Scalability**: Designed to be deployed as a containerized service, the application is built to handle high loads and can be scaled horizontally in a cloud environment.

- **🛡️ Security**: Access to certain API endpoints can be restricted to internal services using middleware that checks for a specific secret value, thereby enhancing the security of the application.

- **🔎 Logging and Monitoring**: Structured logging is implemented using the `zap` library, which facilitates real-time monitoring and debugging of the application.

- **🪫 Graceful Shutdown**: The service can handle shutdown signals and terminate gracefully, ensuring that ongoing operations are completed and resources are released properly.

## Environment Configuration

The following table lists the environment variables used by the Go URL Shortener, which you can set to configure the application:

| Environment Variable    | Description                                                  | Required | Default Value |
|-------------------------|--------------------------------------------------------------|:--------:|:-------------:|
| `DATASTORE_PROJECT_ID`  | Your Google Cloud Datastore project ID.                      | Yes      | None          |
| `INTERNAL_SECRET_VALUE` | A secret value used for internal authentication purposes.    | Yes      | None          |
| `GIN_MODE`              | The mode Gin runs in. Set to "release" for production.       | No       | "debug"       |
| `CUSTOM_BASE_PATH`      | The base path for the URL shortener API endpoints.           | No       | "/"           |

### Notes on Environment Variables

- `DATASTORE_PROJECT_ID` and `INTERNAL_SECRET_VALUE` are mandatory for the application to function correctly. Without these, the application will not be able to connect to Google Cloud Datastore or secure its endpoints.
- `GIN_MODE` is optional and controls the framework's runtime mode. The default mode is "debug", which is suitable for development since it provides detailed logging and error messages. However, it is recommended to set `GIN_MODE` to "release" in a production environment. This turns off debug logging, which can improve performance and prevent the exposure of sensitive information in logs.
- `CUSTOM_BASE_PATH` is optional and allows you to specify a custom base path for all API endpoints. For example, setting this to `/api/v1/` will prefix the routes for retrieving and creating shortened URLs with `/api/v1/`. If not set, the application will use `/` as the default base path.
- Always ensure that environment variables containing sensitive information are kept secure. Do not hardcode them in your application or Dockerfile. Instead, use secure methods of configuration like environment variable injection at runtime or secrets management services.

Remember to set these environment variables before running the application, either locally or as part of your deployment process.

## Getting Started with Docker

<p align="center">
  <img src="https://i.imgur.com/Z1XyTCo.png" width="705" alt="Gopher in Space">
</p>

#### Additional Note:
> The Go URL Shortener can be deployed on Google Cloud Run, which is Kubernetes-compatible, by configuring the necessary environment variables. The Docker image required for deployment is available on GitHub.

The Go URL Shortener can be easily run as a Docker container. Make sure you have Docker installed on your system.

To get the Docker image of Go URL Shortener, pull the image from the GitHub Container Registry:

```sh
docker pull ghcr.io/h0llyw00dzz/go-urlshortener:latest
```

Once you have the image, you can run it as a container:

```sh
docker run -d -p 8080:8080 \
  -e DATASTORE_PROJECT_ID='your-datastore-project-id' \
  -e INTERNAL_SECRET_VALUE='your-internal-secret' \
  ghcr.io/h0llyw00dzz/go-urlshortener:latest
```

This command will start the URL Shortener service and bind it to port 8080 on your host machine.

Make sure to replace `your-datastore-project-id` and `your-internal-secret` with the actual values you want to use for your deployment. These environment variables will be read by your Go application inside the Docker container to configure the connection to Google Cloud Datastore and to set the internal secret for authentication.

### Example Creating a Short URL

To create a short URL, send a `POST` request with a JSON payload containing the original URL. You'll also need to include a custom internal secret header for authentication purposes.

```sh
curl -X POST \
  https://example-your-deployurl-go-dev.a.run.app/ \
  -H 'Content-Type: application/json' \
  -H 'X-Internal-Secret: YOURKEY-SECRET' \
  -d '{"url": "https://go.dev/"}'
```

Replace `YOURKEY-SECRET` with the actual secret key required by your deployment.

The service will respond with a JSON object that includes the ID of the shortened URL:

```json
{
  "id": "{ShortenedID}",
  "shortened_url": "https://example-your-deployurl-go-dev.a.run.app/{ShortenedID}"
}
```

### Example Editing a Short URL

To edit an existing short URL, you will send a `PUT` request with a JSON payload that contains the `id` of the short URL you want to update, the `old_url` which is the current URL associated with that `id`, and the `new_url` that you want to change it to. This operation also requires the custom internal secret header for authentication purposes.

```sh
curl -X PUT \
  https://example-your-deployurl-go-dev.a.run.app/{ShortenedID} \
  -H 'Content-Type: application/json' \
  -H 'X-Internal-Secret: YOURKEY-SECRET' \
  -d '{"id": "{ShortenedID}", "old_url": "https://golang.org/", "new_url": "https://go.dev/"}'
```

Replace `{ShortenedID}` with the actual ID of the shortened URL and `YOURKEY-SECRET` with the actual secret key required by your deployment.

You can then access the shortened URL at `https://example-your-deployurl-go-dev.a.run.app/{ShortenedID}`, which will redirect you to the original URL.

### Example Deleting a Short URL

To delete an existing short URL, you will send a `DELETE` request with a JSON payload that contains both the `id` of the short URL and the `url` that is currently associated with that `id`. This operation may also require authentication, which should be provided through a custom internal secret header.

```sh
curl -X DELETE \
  https://example-your-deployurl-go-dev.a.run.app/{ShortenedID} \
  -H 'Content-Type: application/json' \
  -H 'X-Internal-Secret: YOURKEY-SECRET' \
  -d '{"id": "{ShortenedID}", "url": "https://golang.org/"}'
```

Replace `{ShortenedID}` with the actual ID of the shortened URL you wish to delete, `https://golang.org/` with the actual URL associated with that ID, and `YOURKEY-SECRET` with the actual secret key required by your service for authentication.

## Roadmap

As the project is written in Go, we are considering the development of our own NoSQL database for persistent storage. This would allow us to tailor the storage solution specifically to our needs and avoid dependency on third-party cloud services.

## Contributing

Contributions are very much welcome! Whether you have ideas for improvements or want to help with existing issues, we value your input and collaboration. If you're considering making significant changes or have suggestions, please first open an issue to discuss your ideas. This ensures that we can agree on the best approach and that your efforts align well with the project's direction. To get started and understand our contribution process, please check out our [contributing guidelines](CONTRIBUTING.md). Following these guidelines will help ensure a smooth contribution process and increase the likelihood of your contributions being integrated into the project.
