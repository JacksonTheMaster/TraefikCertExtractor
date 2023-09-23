import json
import os
import time

input_file = '/acme/acme.json'
output_dir = '/extracted-certs'

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
    print("Starting Traefik Acme to x509 Cert Export...")
    extract_certs()
    print("Sleeping 8 minutes, then exporting again...")
    time.sleep(8 * 60)