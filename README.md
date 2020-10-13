# Edge-cloud data pipeline for object detection
- ImageCapture, NATS, YOLO, InfluxDB, and Orchestrator are runs as Docker container microservices
- Images are captured from video camera and published to NATS (Python, ImageCapture microservice)
- The Orchestrator subscribes to images from NATs, and sends to YOLO (Golang, Orchestrator microservice)
- The object detections from YOLO, are then sent to InfluxDB 
- Data from InfluxDB is then sent to AWS IoT Core
- From AWS IoT Core data can then be sent to numerous AWS Services including AWS Timestream table, IoT Analytics for further processing
