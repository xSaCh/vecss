PORT=9090

s3proxy: 
	docker run -d -p 443:443 -p $(PORT):80 --env S3PROXY_AUTHORIZATION=none --env LOG_LEVEL=debug  --env S3PROXY_CORS_ALLOW_ALL=true --env S3PROXY_IGNORE_UNKNOWN_HEADERS=true andrewgaul/s3proxy

s3mock:
	docker run -d -p $(PORT):9090 -p 9191:9191 -e initialBuckets=bkt -t adobe/s3mock

.PHONY: s3proxy s3mock

