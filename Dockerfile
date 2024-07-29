FROM golang:1.22-alpine

# Install CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon@latest

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENTRYPOINT CompileDaemon --build="go build -o main server.go" --command="./main"

EXPOSE 8080
