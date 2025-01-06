   # Use the official Golang image as a build stage
   FROM golang:1.23.4 AS builder

   # Set the Current Working Directory inside the container
   WORKDIR /app

   # Copy go.mod and go.sum files
   COPY go.mod go.sum ./

   # Download all dependencies
   RUN go mod download

   # Copy the entire source code into the container
   COPY . .

   # Change to the directory where main.go is located
   WORKDIR /app/cmd/api

   # Build the Go app
   RUN CGO_ENABLED=0 GOOS=linux go build -o myapp .

   # Start a new stage from scratch
   FROM alpine:latest

   RUN apk add --no-cache postgresql-client

   # Set the Current Working Directory inside the container
   WORKDIR /root/

   # Copy the Pre-built binary file from the previous stage
   COPY --from=builder /app/cmd/api/myapp .

   COPY .env .env
   # Expose port 8080 to the outside world
   EXPOSE 8080

   # Command to run the executable
   CMD ["./myapp"]