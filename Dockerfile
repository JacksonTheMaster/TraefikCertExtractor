FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy source code
COPY main.go /app/

# Build for multiple architectures
RUN GOOS=linux GOARCH=amd64 go build -o traefik-cert-extractor.amd64 main.go && \
    GOOS=linux GOARCH=arm64 go build -o traefik-cert-extractor.arm64 main.go

FROM alpine:3.19

WORKDIR /app

# Copy the compiled binaries from the builder stage
COPY --from=builder /app/traefik-cert-extractor.amd64 /app/
COPY --from=builder /app/traefik-cert-extractor.arm64 /app/

# Copy HTML template
COPY index.html /app/

# Create volumes for acme.json and extracted certificates
VOLUME ["/acme", "/extracted-certs"]

# Run the appropriate binary based on architecture
CMD ["sh", "-c", "if [ \"$(uname -m)\" = \"aarch64\" ]; then exec /app/traefik-cert-extractor.arm64; else exec /app/traefik-cert-extractor.amd64; fi"]