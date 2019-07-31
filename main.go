package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

}
