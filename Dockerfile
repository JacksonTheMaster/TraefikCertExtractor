FROM python:3.9-slim

WORKDIR /app

# Copy scripts
COPY extract_certs.py /app/extract_certs.py

# Copy both binaries
COPY tce.amd64 /app/tce.amd64
COPY tce.arm64 /app/tce.arm64

VOLUME ["/acme", "/extracted-certs"]

# Detect architecture properly
CMD ["sh", "-c", "python -u /app/extract_certs.py && exec /app/tce.$(ARCH)"]
