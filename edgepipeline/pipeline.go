// Edge vision pipeline - Publisher -> NATS -> YOLO -> InfluxDB
// Assumes following services are running - NATS at port 4222, YOLO at port 8080, InfluxDB at port 8086 
// Assumes a database "objdb" already exists. Create using CURL
package edgepipeline

import (
	"log"
	"os"
	"os/signal"

)

const IM_DIR="../test_images"
const NATS_ENDPOINT="nats://172.17.0.4:4222"
const YOLO_ENDPOINT="http://172.17.0.3:8080/detect"
const DB_ENDPOINT="http://172.17.0.2:8086"
const db="objdb"
const TOPIC="images"

func Pipeline() {


	imCh := make(chan []byte)
	defer close(imCh)

	objCh := make(chan string)
	defer close(objCh)

	// Cleanup on Ctrl-c
	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)

	signal.Notify(sigCh, os.Interrupt, os.Kill)

	go subscribe(NATS_ENDPOINT, TOPIC, imCh)
	go objDetect(YOLO_ENDPOINT, imCh, objCh)
	go writeDB(DB_ENDPOINT, db, objCh)

	// Exit with Ctrl-C 
	<-sigCh
	log.Println("Exiting...")

}

