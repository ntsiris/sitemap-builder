FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

# Create a smaller runtime image
FROM alpine:3.18
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/bin/sitemap /app/sitemap

# Set the entry point to the sitemap binary
ENTRYPOINT ["/app/sitemap"]

# Default command with options for URL and depth (can be overridden)
CMD ["-url=https://example.com", "-depth=3", "-out=/app/map.xml"]
