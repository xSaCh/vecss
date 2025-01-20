import requests
import upload as upld
import json

FILE_PATH = './output.mp4'
upload_url = "http://localhost:8080/upload/"
combine_url = "http://localhost:8080/combine/"

files = { "file": f"open({FILE_PATH}, 'rb')" }

response = requests.post(upload_url, files=files)

if response.status_code != 200:
    print(f"Error Upload with {response.status_code} status: {response.reason} text: {response.text}")

res = response.json()

if 'error' in res:
    print(f"Error in Upload: {res}")
    exit(1)

with open("presigned.json",'w') as f:
    f.write(response.text)

try:
    upld_res = upld.upload(FILE_PATH)
except Exception as e:
    print(e)
    exit(1)

print(upld_res['part_numbers'])
payload = upld_res
headers = {"content-type": "application/json"}

cbn_response = requests.post(combine_url, json=payload, headers=headers)

if response.status_code != 200:
    print(f"Error Combine with {cbn_response.status_code} status: {cbn_response.reason} text: {cbn_response.text}")

cbn_res = cbn_response.json()

if 'error' in cbn_res:
    print(f"Error in Combine: {cbn_res}")
    exit(1)

print(cbn_res)


