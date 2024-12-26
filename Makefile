PORT=9090

s3: 
	docker run -d -p 443:443 -p $(PORT):80 --env S3PROXY_AUTHORIZATION=none --env LOG_LEVEL=debug  --env S3PROXY_CORS_ALLOW_ALL=true --env S3PROXY_IGNORE_UNKNOWN_HEADERS=true andrewgaul/s3proxy

.PHONY: s3

