import requests
import progressbar
from pathlib import Path

def main():
	data = eval(open('presigned.json').read())
	upload_by = data['upload_by']
	max_size = data['max_size']
	urls = data['urls']
	target_file = Path(data['file_name'])
	file_size = data['file_size']
	key = data['key']
	upload_id = data['upload_id']
	bucket_name = data['bucket_name']
	bar = progressbar.ProgressBar(maxval=file_size, \
    	widgets=[progressbar.Bar('=', '[', ']'), ' ', progressbar.Percentage()])
	json_object = dict()
	parts = []
	file_size_counter = 0
	with target_file.open('rb') as fin:
		bar.start()
		for num, url in enumerate(urls):
			part = num + 1
			file_data = fin.read(max_size)
			file_size_counter += len(file_data)
			res = requests.put(url, data=file_data)
			
			if res.status_code != 200:
				print(res.status_code)
				print("Error while uploading your data.")
				return None
			bar.update(file_size_counter)
			etag = res.headers['ETag']
			parts.append((etag, part))
		bar.finish()
		json_object['parts'] = [{"ETag": eval(x), 'PartNumber': int(y)} for x, y in parts]
		json_object['upload_id'] = upload_id
		json_object['key'] = key
		json_object['bucket_name'] = bucket_name
	# requests.post('https://YOUR_HOSTED_API/combine', json={'parts': json_object})
	print(json_object)
	print("Dataset is uploaded successfully")

if __name__ == "__main__":
	main()
