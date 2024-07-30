# Start from Golang base image
FROM golang:alpine as builder

# Enable go modules
ENV GO111MODULE=on

# Install git
RUN apk update && apk add --no-cache git

# Set working directory
WORKDIR /todo

# Copy and download dependencies
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy source code and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/main .

# Start a new stage from scratch
FROM scratch

# Copy binary and .env file
COPY --from=builder /todo/bin/main .
COPY --from=builder /todo/env .env

# Run executable
CMD ["./main"]
