import json
import requests
import progressbar
from pathlib import Path
import os
import sys

def upload(file_path):
	data = eval(open('presigned.json').read())
	# max_size = 5*1024*1024
	max_size = data['chunk_size']
	urls = data['urls']
	target_file = Path(file_path)
	file_size = os.path.getsize(target_file)

	bar = progressbar.ProgressBar(maxval=file_size, \
    	widgets=[progressbar.Bar('=', '[', ']'), ' ', progressbar.Percentage()])
	json_object = dict()
	json_object['upload_id'] = data['upload_id']
	json_object['key'] = target_file.name

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
		# json_object['parts'] = [{"ETag": eval(x), 'PartNumber': int(y)} for x, y in parts]
		json_object['etags'] = [eval(x) for x,y in parts]
		json_object['part_numbers'] = [int(y) for x,y in parts]
	return json_object

if __name__ == "__main__":
	print(upload(sys.argv[1]))
