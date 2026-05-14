# Use the official Go image (version aligned with go.mod)
# https://hub.docker.com/_/golang
FROM golang:1.25-bookworm

# Create and change to the app directory.
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Copy local code to the container image.
COPY . ./

# Install project dependencies
RUN go get .

RUN apt-get update && apt-get install -y gcc

# Build the app
RUN CGO_ENABLED=1 GOOS=linux go build -o app

# Run the service on container startup.
ENTRYPOINT ["./app"]
