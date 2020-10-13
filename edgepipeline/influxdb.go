// Edge vision pipeline - InfluxDB for data storage
package edgepipeline

import (
	"strings"
	"github.com/influxdata/influxdb1-client/v2"
	"log"
	"time"
	"strconv"
)


// Write to Influx DB 
// Key Fields, value - objects detected, probability
// Tag fields, objects class, object detected



func writeDB(dbEndpoint, db string,  objCh <-chan string) {
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: dbEndpoint,
	})
	if err != nil {
		log.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "s",
	})


	for obj := range objCh {
		var fKey string
		var fVal float64
		objects := make(map[string][]float64)
		tags := make(map[string]string)
		fields := make(map[string]interface{})

		// Extract detected object and probability
		for _, ele := range strings.Split(obj, "],") {
				for i, val := range strings.Split(ele, "\n"){
				if i == 2 {
					val = strings.TrimSpace(val)
					fKey = strings.Trim(val, "\",")
				}
				if i == 3 {
					val = strings.Trim(val, ",")
					val = strings.TrimSpace(val)
					fVal, _ = strconv.ParseFloat(val, 64)
				}
			}
			objects[fKey] = append(objects[fKey], fVal)
			for k,v  := range objects {
				tags["source-id"] = "cam1"
				fields["object"] = k
				fields["prob"] = v
			}
			pt, err := client.NewPoint("objdata", tags, fields, time.Now())
			if err != nil {
					log.Println("InfluxDB Error: ", err.Error())
			}
			bp.AddPoint(pt)
		}

		c.Write(bp)

	}
}


