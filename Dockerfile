# When running continaer use flag --privileged --device=/dev/video0:/dev/video0 to access webcam
FROM python:3

# Move to working directory /build
WORKDIR /build

# Copy the code into the container
COPY frame_capture_nats.py .

# Numpy
RUN pip install numpy

# OpenCV
RUN pip install opencv-python

# NATS
RUN pip install nats-python

# Execute
CMD [ "python", "/build/frame_capture_nats.py" ]

# For Debug using attach shell
#ENTRYPOINT ["tail", "-f", "/dev/null"]
