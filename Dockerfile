# Start with a minimal Go image
FROM golang:1.21-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o tv-bot ./cmd/bot

# Expose any necessary ports
# EXPOSE 8080

# Run the executable
CMD ["./tv-bot"]
