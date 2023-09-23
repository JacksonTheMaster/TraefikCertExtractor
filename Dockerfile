FROM python:3.9-slim

WORKDIR /app

COPY extract_certs.py /app/extract_certs.py

RUN pip install watchdog

VOLUME ["/acme.json", "/extracted-certs"]

CMD ["python", "/app/extract_certs.py"]
