# Use the official Golang image as the base image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum file into the container
COPY go.mod .

# Download the dependencies
RUN go mod download

# Copy the server code into the container
COPY server.go .

# Build the server binary
RUN go build -o server .

# Set the command to run the server binary when the container starts
ENTRYPOINT ["/app/server"]