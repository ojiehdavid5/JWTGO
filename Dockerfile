# Use the official Golang image
FROM golang:1.20 AS builder

# Set the working directory
WORKDIR /app
# Copy the Go modules manifests
COPY ./go.mod ./go.sum ./
# Download the Go modules
RUN go mod download
# Copy the source code
COPY . .
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .
# Use a minimal base image for the final build
FROM alpine:latest
# Set the working directory
WORKDIR /root/
# Copy the binary from the builder stage
COPY --from=builder /app/main .
# Expose the port the app runs on
EXPOSE 8080
# Command to run the binary
CMD ["./main"]
# Healthcheck to ensure the app is running  
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/ || exit 1
# Metadata
LABEL org.opencontainers.image.title="Go Web App" \
      org.opencontainers.image.description="A simple Go web application" \
      org.opencontainers.image.version="1.0.0" \
      org.opencontainers.image.authors="Your Name <ojiehdavid5@gmail.com>" \
      org.opencontainers.image.url="https://github.com/chuks/JWTGO"
# Add a label for the license
LABEL org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.source="https://github.com/chuks/JWTGO"



