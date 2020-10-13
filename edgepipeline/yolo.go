// Edge vision pipeline - YOLO 
package edgepipeline

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"

)



func objDetect(yoloEndpoint string, imChan <-chan []byte, objChan chan<- string) {
// Infinite loop until channel is closed
	for im := range imChan {
		imbuf := bytes.NewBuffer(im)
		log.Println("YOLO objDetect", len(im))

		// Prepare a form to submit to YOLO URL
		var b bytes.Buffer
		w := multipart.NewWriter(&b)
		fw, err := w.CreateFormFile("image_file", "imfile")
		if err != nil {
			log.Fatal(err)
		}
		if _, err = io.Copy(fw, imbuf); err != nil {
			log.Fatal(err)
		}

		// Add the other fields - detection threshold
		if fw, err = w.CreateFormField("threshold"); err != nil {
			log.Fatal(err)
		}
		if _, err = fw.Write([]byte("0.25")); err != nil {
			log.Fatal(err)
		}
		// Close the multipart writer.
		w.Close()

		// Submit it to handler.
		req, err := http.NewRequest("POST", yoloEndpoint, &b)
		if err != nil {
			log.Fatal(err)
		}

		// Set the content type, this will contain the boundary.
		req.Header.Set("Content-Type", w.FormDataContentType())
		// Submit the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// Send the response on output channel
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			log.Println("YOLO object description", bodyString)
			objChan <- bodyString
		}

	}
}
