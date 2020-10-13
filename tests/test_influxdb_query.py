# Query to test write into InfluxDB 
from influxdb import InfluxDBClient


# Query to verify insert into InfluxDB 
influxDBHost = '127.0.0.1'
dbname = "objdb"

# Connect to influx dB
client = InfluxDBClient(host=influxDBHost, port=8086)
client.create_database(dbname)
client.switch_database(dbname)

def main():
	q1 = "select * from objdata;"
	queryData(q1)
	#q2 = "select prob from objdata;"
	#queryData(q2)	

# Query data
def queryData(query):
	print("Querying data: " + query)
	result = client.query(query)
	#print("Result: {0}".format(result))
	eventData = list(result.get_points(measurement='objdata'))
	print(eventData)
	for ele in eventData:
		for k, v in ele.items():
			print(k, v)

if __name__ == "__main__":
    main()

