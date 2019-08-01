package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"net/http"

	"github.com/nfnt/resize"
)

type Message struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

func main() {

	http.HandleFunc("/test", HomeHandler)
	http.HandleFunc("/image", ImageHandler)

	fmt.Println("Server is started...")
	http.ListenAndServe(":8989", nil)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//response := fmt.Sprintf("I'm your http handler")
	fmt.Fprintf(w, "I'm your http handler")
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	output, err := json.Marshal(msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)

	sourceImage, err := base64.StdEncoding.DecodeString(msg.Body)
	fmt.Println("body:.............................: " + msg.Body)

	img1, _, err := image.Decode(bytes.NewReader(sourceImage))

	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return
	}

	m := resize.Resize(200, 200, img1, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, m, nil)
	imageBit := buf.Bytes()
	/*Defining the new image size*/

	photoBase64 := base64.StdEncoding.EncodeToString([]byte(imageBit))
	fmt.Println("Photo Base64.............................: " + photoBase64)

}
