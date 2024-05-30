# Use the official Golang image as the base image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod file into the container
COPY go.mod .

# Download the dependencies
RUN go mod download

# Copy the client code into the container
COPY client.go .

# Build the client binary
RUN go build -o client .

# Set the entrypoint command
ENTRYPOINT ["/app/client"]
