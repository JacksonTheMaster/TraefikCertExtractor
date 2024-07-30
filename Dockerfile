FROM python:3.9-slim

WORKDIR /app

COPY extract_certs.py /app/extract_certs.py
COPY tce.amd64 /app/tce.amd64

VOLUME ["/acme", "/extracted-certs"]

CMD ["python", "-u", "/app/extract_certs.py"]
