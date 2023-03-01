package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello This is the home"))
}
func view(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the snippet"))
}
func create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a snippet"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", view)
	mux.HandleFunc("/snippet/create", create)
	fmt.Println("Server Starting on port 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
