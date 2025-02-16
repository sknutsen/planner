# Use the Go 1.23 alpine official image
# https://hub.docker.com/_/golang
FROM golang:1.23.4-alpine

# Create and change to the app directory.
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy local code to the container image.
COPY . ./

# Install project dependencies
RUN go get .

RUN apk add build-base

# Build the app
RUN CGO_ENABLED=1 go build -o app

# Run the service on container startup.
ENTRYPOINT ["./app"]
