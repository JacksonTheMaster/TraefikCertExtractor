import json
import os
import time
import base64
import subprocess
import platform

input_file = '/acme/acme.json'
output_dir = '/extracted-certs'

# Correct architecture mapping
arch_map = {
    "x86_64": "amd64",
    "aarch64": "arm64"
}

arch = platform.machine()
executable_path = f"./tce.{arch_map.get(arch, 'amd64')}"  # Default to amd64 if unknown

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
            decoded_certificate = base64.b64decode(certificate).decode('utf-8')
            decoded_key = base64.b64decode(key).decode('utf-8')
        except Exception as e:
            print(f"Error decoding the base64 data for domain {domain}. Reason: {str(e)}")
            continue

        try:
            with open(os.path.join(output_dir, f'{domain}.crt'), 'w') as cert_file:
                cert_file.write(decoded_certificate)
        except Exception as e:
            print(f"Error: Unable to write certificate for domain {domain}. Reason: {str(e)}")
            continue

        try:
            with open(os.path.join(output_dir, f'{domain}.key'), 'w') as key_file:
                key_file.write(decoded_key)
        except Exception as e:
            print(f"Error: Unable to write key for domain {domain}. Reason: {str(e)}")
            continue

if __name__ == "__main__":
    print("Starting Traefik Acme to x509 Cert Export...")
    extract_certs()
    print(f"Detected architecture: {arch}")
    print(f"Using binary: {executable_path}")
    print("Sleeping 8 minutes, then exporting again...")

    try:
        subprocess.run([executable_path], check=True)
    except subprocess.CalledProcessError as e:
        print(f"Error running {executable_path}. Reason: {str(e)}")
    except OSError as e:
        print(f"Execution failed: {executable_path} is not a valid binary for this architecture. Reason: {str(e)}")
    
    time.sleep(8 * 60)
