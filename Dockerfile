# Use the official Golang image to create a build artifact.
FROM golang:1.21.5 as builder

# Copy local code to the container image.
WORKDIR /app
COPY go.mod go.sum ./

# Verify the integrity of the modules.
RUN go mod verify

# Copy the rest of the code
COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o go-urlshortner

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Add a non-root user called 'go-urlshortner'
RUN adduser -D go-urlshortner

# Copy the built binary from the builder stage.
COPY --from=builder /app/go-urlshortner /go-urlshortner

# Use the non-root user to run our application
USER go-urlshortner

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the web service on container startup.
CMD ["/go-urlshortner"]

# Example run with ENV variables:
# docker build -t go-urlshortner .
# docker run -p 8080:8080 -e DATASTORE_PROJECT_ID='your-project-id' go-urlshortner
# Note: Support deploy in k8s with env variables
