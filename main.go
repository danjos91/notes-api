// Notes API - Step 3
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Note struct {
	ID      int
	Title   string
	Content string
}

var store = map[string]Note{}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	//Check method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//Initialize notes to return
	var notes []Note
	//Loop to get the notes from store
	for _, note := range store {
		notes = append(notes, note)
	}
	//Getting ready de json to return, data is binary []byte
	data, err := json.Marshal(notes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Set the content type to json
	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK) //200 by default
	w.Write(data)
}

func main() {
	store["1"] = Note{ID: 1, Title: "First Note", Content: "Just a note"}
	store["2"] = Note{ID: 2, Title: "Second Note", Content: "It's getting interesting"}
	store["3"] = Note{ID: 3, Title: "Third Note", Content: "Muy interesante!"}
	store["4"] = Note{ID: 4, Title: "Fourth Note", Content: "Hola Amigo! Como estas?"}
	fmt.Println("Notes API")

	fmt.Println("Server is running on port 8080")

	http.HandleFunc("/notes", notesHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
