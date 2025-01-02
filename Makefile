PORT=9090

run_s3proxy: 
	docker run -d -p 443:443 -p $(PORT):80 --name vecss_s3proxy --env S3PROXY_AUTHORIZATION=none --env LOG_LEVEL=debug  --env S3PROXY_CORS_ALLOW_ALL=true --env S3PROXY_IGNORE_UNKNOWN_HEADERS=true andrewgaul/s3proxy

run_s3mock:
	docker kill vecss_s3mock
	docker rm vecss_s3mock
	docker run -d -p $(PORT):9090 -p 9191:9191 -e initialBuckets=bkt --name vecss_s3mock -t adobe/s3mock

start_s3mock:
	docker start vecss_s3mock

.PHONY: run_s3proxy run_s3mock start_s3mock

