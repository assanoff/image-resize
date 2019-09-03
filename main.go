package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/nfnt/resize"
)

type Message struct {
	ID     string `json:"id"`
	Size   uint   `json: "size"`
	Body   string `json:"body"`
	Format string `json:"format"`
}

func main() {

	http.HandleFunc("/test", HomeHandler)
	http.HandleFunc("/image", ImageHandler)

	fmt.Println("Server is started...")
	http.ListenAndServe(":3220", nil)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//response := fmt.Sprintf("I'm your http handler")
	fmt.Fprintf(w, "I'm your http handler")
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	//Read body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//Parse JSON
	var msg Message
	err = json.Unmarshal(b, &msg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// output, err := json.Marshal(msg)
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	// w.Header().Set("content-type", "application/json")
	// w.Write(output)

	//Decode string to image
	sourceImage, err := base64.StdEncoding.DecodeString(msg.Body)
	//fmt.Println("body:.............................: " + msg.Body)

	//Decode from bytes to image
	imgSource, _, err := image.Decode(bytes.NewReader(sourceImage))

	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return
	}

	//Resize image
	m := resize.Resize(msg.Size, msg.Size, imgSource, resize.Lanczos3)
	buf := new(bytes.Buffer)
	format := strings.ToLower(msg.Format)

	switch format {
	case "jpg":
		err = jpeg.Encode(buf, m, nil)
	case "jpeg":
		err = jpeg.Encode(buf, m, nil)
	case "png":
		err = png.Encode(buf, m)
	default:
		fmt.Println("Unknown format")
	}

	imageBit := buf.Bytes()
	/*Defining the new image size*/

	photoBase64 := base64.StdEncoding.EncodeToString([]byte(imageBit))
	//fmt.Println("Photo Base64.............................: " + photoBase64)

	outputMsg := Message{}
	outputMsg.ID = msg.ID
	outputMsg.Size = msg.Size
	outputMsg.Body = photoBase64

	output, err := json.Marshal(outputMsg)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.Write(output)
	fmt.Printf("%s\t%s\n", outputMsg.ID, format)
}
