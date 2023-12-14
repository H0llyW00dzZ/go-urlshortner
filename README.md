# Go URL Shortener

[![Go Report Card](https://goreportcard.com/badge/github.com/H0llyW00dzZ/go-urlshortner)](https://goreportcard.com/report/github.com/H0llyW00dzZ/go-urlshortner)
[![Go Reference](https://pkg.go.dev/badge/github.com/H0llyW00dzZ/go-urlshortner.svg)](https://pkg.go.dev/github.com/H0llyW00dzZ/go-urlshortner)

This is a simple URL shortening service written in Go, which utilizes the Gin Web Framework for routing and Google Cloud Datastore for persistent storage.

## Introduction

<p align="center">
  <img src="https://i.imgur.com/T04JNPd.jpg" alt="Go Picture">
</p>

This project aims to provide a straightforward and scalable approach to creating short aliases for lengthy URLs. It's constructed in Go and is designed to be simple to deploy and maintain. The service includes basic functionalities such as generating a shortened URL and redirecting to the original URL when accessed.

## Features

- **Shorten URLs**: Convert long URLs into short, manageable links that are easier to share.
- **Redirection**: Use the generated short link to redirect to the original URL.
- **Simple Integration**: Easily integrate with your applications using RESTful API endpoints.

## Getting Started with Docker

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

You can then access the shortened URL at `https://example-your-deployurl-go-dev.a.run.app/{ShortenedID}`, which will redirect you to the original URL.

## Roadmap

As the project is written in Go, we are considering the development of our own NoSQL database for persistent storage. This would allow us to tailor the storage solution specifically to our needs and avoid dependency on third-party cloud services.

## Contributing

Contributions are very much welcome! If you're thinking about making significant changes or improvements, please start by opening an issue. This allows us to have a discussion about the proposed changes and agree on the best approach before moving forward. We value your ideas and input, and we look forward to collaborating with you!
