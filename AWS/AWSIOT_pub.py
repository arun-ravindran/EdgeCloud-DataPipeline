# Adapted from pubsub.py -  https://github.com/aws/aws-iot-device-sdk-python-v2
# Publishes data read from InfluxDB to IoT core using MQTT 

from __future__ import absolute_import
from __future__ import print_function
import argparse
from awscrt import io, mqtt, auth, http
from awsiot import mqtt_connection_builder
import sys
import time

import dbRead # For reading from InfluxDB

parser = argparse.ArgumentParser(description="Send and receive messages through and MQTT connection.")
parser.add_argument('--endpoint', required=True, help="Your AWS IoT custom endpoint, not including a port. " +
                                                      "Ex: \"abcd123456wxyz-ats.iot.us-east-1.amazonaws.com\" See IoT settings")
parser.add_argument('--cert', help="File path to your client certificate, in PEM format.")
parser.add_argument('--key', help="File path to your private key, in PEM format.")
parser.add_argument('--root-ca', help="File path to root certificate authority, in PEM format. " +
                                      "Necessary if MQTT server uses a certificate that's not already in " +
                                      "your trust store.")
parser.add_argument('--client-id', default='edgeserver', help="Client ID for MQTT connection.")
parser.add_argument('--topic', default="edgeserver/objdata", help="Topic to subscribe to, and publish messages to.")

# Using globals to simplify sample code
args = parser.parse_args()


# Callback when connection is accidentally lost.
def on_connection_interrupted(connection, error, **kwargs):
    print("Connection interrupted. error: {}".format(error))


# Callback when an interrupted connection is re-established.
def on_connection_resumed(connection, return_code, session_present, **kwargs):
    print("Connection resumed. return_code: {} session_present: {}".format(return_code, session_present))

if __name__ == '__main__':
	# Spin up resources
	event_loop_group = io.EventLoopGroup(1)
	host_resolver = io.DefaultHostResolver(event_loop_group)
	client_bootstrap = io.ClientBootstrap(event_loop_group, host_resolver)

	mqtt_connection = mqtt_connection_builder.mtls_from_path(
		endpoint=args.endpoint,
		cert_filepath=args.cert,
		pri_key_filepath=args.key,
		client_bootstrap=client_bootstrap,
		ca_filepath=args.root_ca,
		on_connection_interrupted=on_connection_interrupted,
		on_connection_resumed=on_connection_resumed,
		client_id=args.client_id,
		clean_session=False,
		keep_alive_secs=6)

	print("Connecting to {} with client ID '{}'...".format(
		args.endpoint, args.client_id))

	connect_future = mqtt_connection.connect()

	# Future.result() waits until a result is available
	connect_future.result()
	print("Connected!")


    # Publish messageis read from InfluxDB
	q1 = "select * from objdata;" # Query
	messages = dbRead.queryData(q1)
	publish_count = 1
	for msg in messages:
		print("Publishing message to topic '{}': {}".format(args.topic, msg))
		mqtt_connection.publish(
			topic=args.topic,
			payload=msg,
			qos=mqtt.QoS.AT_LEAST_ONCE)
		publish_count += 1
		time.sleep(1)
	print("{} number of messages published".format(publish_count))

	# Disconnect
	print("Disconnecting...")
	disconnect_future = mqtt_connection.disconnect()
	disconnect_future.result()
	print("Disconnected!")
