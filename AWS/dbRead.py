# Query InfluxDB and return a list of json responses
from influxdb import InfluxDBClient
import json

# Query to verify insert into InfluxDB 
influxDBHost = '127.0.0.1'
dbname = "objdb"

# Connect to influx dB
client = InfluxDBClient(host=influxDBHost, port=8086)
client.create_database(dbname)
client.switch_database(dbname)

def main():
	q1 = "select * from objdata;"
	print(queryData(q1))
	#q2 = "select prob from objdata;"
	#queryData(q2)	

# Query data
messages = []
def queryData(query):
	result = client.query(query)
	eventData = list(result.get_points(measurement='objdata'))
	print(eventData)
	for ele in eventData:
		msg = json.dumps(ele) # Convert to JSON
		messages.append(msg)
	return messages
	
if __name__ == "__main__":
    main()

