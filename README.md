# VECSS

**Video Transpiling Streaming System (VECSS)** demonstrates a distributed system for efficient video processing by utilizing modular servers. The system is organized as follows:

- **VUS (Upload Server):**
  - Handles all upload-related operations.
  - Uploads videos to an S3 bucket using s3_mock for testing (easily replaceable with any S3-compatible storage).
  - Creates tasks and pushes them to a messaging queue.

- **VTS (Transpiling Worker Server):**
  - Acts as a worker server group that picks up video transcoding tasks from the messaging queue (RabbitMQ).
  - Requires FFMPEG to be installed for processing videos.

- **(IN PROGRESS) VSS (Streaming Server):**
  - Manages all streaming-related functionalities for processed videos.

## Installation

1. **Docker Images:**
   - Pull the Adobe s3_mock and RabbitMQ images:
     ```bash
     docker pull adobe/s3_mock
     docker pull rabbitmq:management
     ```

2. **Start Docker Services:**
   - Run the local S3 server:
     ```bash
     make run_s3mock
     ```
   - Run the RabbitMQ server:
     ```bash
     make run_mq
     ```

3. **Run Servers:**
   - Start each server using:
     ```bash
     go run vus/vus.go
     go run vts/vts.go
     ```

## Usage

- **Testing:**  
  Run the provided test script to verify everything is working:
  ```bash
  python test/test_vus_vts.py
