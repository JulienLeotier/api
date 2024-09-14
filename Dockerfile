FROM golang:1.22-alpine

# Install CompileDaemon for live-reloading
RUN go install github.com/githubnemo/CompileDaemon@latest

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the code to the container
COPY . .

# Use CompileDaemon to build and run the app
ENTRYPOINT CompileDaemon --build="go build -o main cmd/main.go" --command="./main"

# Expose port 8080 for the app
EXPOSE 8080
