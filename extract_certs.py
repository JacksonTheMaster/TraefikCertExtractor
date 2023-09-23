import json
import os
import time
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler

input_file = '/acme.json'
output_dir = '/extracted-certs'

class AcmeJsonHandler(FileSystemEventHandler):
    def on_modified(self, event):
        if not event.is_directory and event.src_path == input_file:
            print("acme.json modified, extracting certificates...")
            extract_certs()

def extract_certs():
    if not os.path.exists(input_file):
        print(f"Error: {input_file} does not exist.")
        return

    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    try:
        with open(input_file, 'r') as f:
            data = json.load(f)
    except json.JSONDecodeError:
        print(f"Error: {input_file} is not a valid JSON file.")
        return
    except Exception as e:
        print(f"Error: Unable to open {input_file}. Reason: {str(e)}")
        return

    for cert_info in data.get('letsencrypt', {}).get('Certificates', []):
        domain = cert_info['domain']['main']
        certificate = cert_info.get('certificate')
        key = cert_info.get('key')

        if not certificate or not key:
            print(f"Error: Certificate or key for domain {domain} is empty.")
            continue

        try:
            with open(os.path.join(output_dir, f'{domain}.crt'), 'w') as cert_file:
                cert_file.write("-----BEGIN CERTIFICATE-----\n")
                cert_file.write(certificate)
                cert_file.write("\n-----END CERTIFICATE-----\n")
        except Exception as e:
            print(f"Error: Unable to write certificate for domain {domain}. Reason: {str(e)}")
            continue

        try:
            with open(os.path.join(output_dir, f'{domain}.key'), 'w') as key_file:
                key_file.write("-----BEGIN PRIVATE KEY-----\n")
                key_file.write(key)
                key_file.write("\n-----END PRIVATE KEY-----\n")
        except Exception as e:
            print(f"Error: Unable to write key for domain {domain}. Reason: {str(e)}")
            continue

if __name__ == "__main__":
    extract_certs()

    observer = Observer()
    observer.schedule(AcmeJsonHandler(), path=os.path.dirname(input_file))
    observer.start()

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        observer.stop()
    observer.join()
