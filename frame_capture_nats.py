# Capture frames from webcam using OpenCV and publish to NATS
# pip3 install nats-python
# NATS server publisher - run subscriber, followed by producer


import cv2
import signal
import sys
import time
from pynats import NATSClient
import os

NATS_URL="nats://172.17.0.4:4222"
TOPIC="images"


def signalHandler(sig, frame):
	print("Exiting...")
	sys.exit(0)

signal.signal(signal.SIGINT, signalHandler)
vidcap = cv2.VideoCapture(0)
success,image = vidcap.read()
if success: print("Opened camera succesfully..Ctrl-C to exit")
count = 1
while success:
	_, imageJPG = cv2.imencode(".jpg", image)
	imageByte = imageJPG.tobytes()  	
	with NATSClient(NATS_URL, socket_timeout=2) as client:
		client.publish(TOPIC, payload=imageByte)
	time.sleep(1)      
	success,image = vidcap.read()
	count += 1
vidcap.release()

print("Captured {} frames".format(count))

