FROM python:3.9-slim

WORKDIR /app

COPY extract_certs.py /app/extract_certs.py

VOLUME ["/acme", "/extracted-certs"]

CMD ["python", "-u", "/app/extract_certs.py"]
