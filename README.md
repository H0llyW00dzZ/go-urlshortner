# Go URL Shortener

[![Go Report Card](https://goreportcard.com/badge/github.com/H0llyW00dzZ/go-urlshortner)](https://goreportcard.com/report/github.com/H0llyW00dzZ/go-urlshortner)

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

## Getting Started

To use the Go URL Shortener, you'll need to make HTTP requests to its API endpoints. Currently, the service supports creating short URLs with a `POST` request.

### Creating a Short URL

To create a short URL, send a `POST` request with a JSON payload containing the original URL. You'll also need to include a custom internal secret header for authentication purposes.

```sh
$ curl -X POST \
  https://example-your-deployurl-go-dev.a.run.app/ \
  -H 'Content-Type: application/json' \
  -H 'X-Internal-Secret: YOURKEY-SECRET' \
  -d '{"url": "https://go.dev/"}'
```

Replace `YOURKEY-SECRET` with the actual secret key required by your deployment.

The service will respond with a JSON object that includes the ID of the shortened URL:

```json
{"id":"ShortenedID"}
```

You can then access the shortened URL at `https://example-your-deployurl-go-dev.a.run.app/{ShortenedID}`, which will redirect you to the original URL.

## Contributing

Contributions are welcome! For major changes, please open an issue first to discuss what you would like to change.
