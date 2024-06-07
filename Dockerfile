# Stage 1: Build the Go application
FROM golang:1.17 as builder

WORKDIR /app

# Copy go mod and sum files
#COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
#RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o memos-webhook

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/memos-webhook .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./memos-webhook"]
