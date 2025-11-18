FROM golang:1.24-alpine

# Environment variable
WORKDIR /usr/src/app

# Install Air for hot-reloading in development (air-verse/air is the new path)
RUN go install github.com/air-verse/air@v1.61.0

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Expose the server on port 3000
EXPOSE 3000

# Air will be used for hot-reloading during development
# For production, use: CMD ["go", "run", "cmd/main.go"]
