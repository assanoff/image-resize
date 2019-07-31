package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/test", HomeHandler)
	fmt.Println("Server is started...")
	http.ListenAndServe(":8989", nil)

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	//response := fmt.Sprintf("I'm your http handler")
	fmt.Fprintf(w, "I'm your http handler")
}
