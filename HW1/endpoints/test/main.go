package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var api_key = "&lt;acc_7d54f0cae7319d2&gt;"
var api_secret = "&lt;d4a4b8a3f75b7539bef24cd181f070d0&gt;"

var jsonString = ` 
{
    "result": {
        "faces": [
            {
                "confidence": 99.99755859375,
                "coordinates": {
                    "height": 122,
                    "width": 122,
                    "xmax": 387,
                    "xmin": 265,
                    "ymax": 156,
                    "ymin": 34
                },
                "face_id": "60577279bdedbcfbd8b4186ad5a9bd94f89a9085985d0edd41ad38298058c44c"
            }
        ]
    },
    "status": {
        "text": "",
        "type": "success"
    }
}`

type Coordinates struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	XMax   int `json:"xmax"`
	XMin   int `json:"xmin"`
	YMax   int `json:"ymax"`
	YMin   int `json:"ymin"`
}

type Face struct {
	Confidence  float64     `json:"confidence"`
	Coordinates Coordinates `json:"coordinates"`
	FaceID      string      `json:"face_id"`
}

type Result struct {
	Faces []Face `json:"faces"`
}

type Status struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type JSONData struct {
	Result Result `json:"result"`
	Status Status `json:"status"`
}

func main() {
	// var data map[string]interface{}
	// _ = json.NewDecoder(http.Response()).Decode(&data)
	var j JSONData
	_ = json.Unmarshal([]byte(jsonString), &j)
	fmt.Printf("status: %v\nfaceid: %v\n ", j.Status.Type, j.Result.Faces[0].FaceID)
	// Your string that you want to convert to a MIME header
	// fileHeader := &multipart.FileHeader{
	// 	Filename: "example.jpg",
	// }
	// fileExtension := filepath.Ext(fileHeader.Filename)
	// fmt.Println(fileExtension)
}

func DetectFace(url string) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "https://api.imagga.com/v2/faces/detections", nil)
	req.SetBasicAuth(api_key, api_secret)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error when sending request to the server")
		return "", err
	}
	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Printf("Couldn't decode response body into map. Here's why: %v\n", err)
		return "", err
	}
	s, ok := data["face_id"].(string)
	if !ok {
		log.Printf("response json: %v\n", data)
		log.Printf("Internal server error: invalid face id format: %v, %T\n", data["face_id"], data["face_id"])
		return "", nil
	}
	return s, nil
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	buf, err := os.Open("image.jpg")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("image", "image.jpg")

	_, err = io.Copy(fw, buf)
	if err != nil {
		log.Fatal(err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", "http://127.0.0.1:80/v1/vision/custom/testmodel", bytes.NewReader(body.Bytes()))
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{
		// Set timeout to not be at mercy of microservice to respond and stall the server
		Timeout: time.Second * 20,
	}

	rsp, _ := client.Do(req)

	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}
	log.Print(rsp.StatusCode)

	body2, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body2))

	w.Header().Set("Content-Type", "JSON")
	w.Write(body2)
}
