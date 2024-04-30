FROM golang:1.21-bookworm as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /compiler-server ./server.go

# Start a new stage from scratch
FROM docker:26.1.0-dind

# Set the Current Working Directory inside the container
WORKDIR /

# Import the user and group files from the builder
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /compiler-server .

# Expose the port the server will run on
EXPOSE 8080

# Command to run the executable
CMD ["/compiler-server"]