// Notes API - Step 3
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Note struct {
	ID      int
	Title   string
	Content string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var store = map[string]Note{}

func notesHandler(w http.ResponseWriter, r *http.Request) {
	//Check method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//return value of type json
	w.Header().Set("Content-Type", "application/json")
	//check if id is present
	id := strings.Trim(r.URL.Path, "/notes/")

	var data []byte
	var err error

	if len(id) > 0 {
		note := store[id]
		if note.ID == 0 { //actually checking if not is not zero
			w.WriteHeader(http.StatusNotFound)
			data, _ = json.Marshal(ErrorResponse{Error: "not found"})
		} else {
			data, err = json.Marshal(store[id])
		}
	} else {
		var notes []Note
		for _, note := range store {
			notes = append(notes, note)
		}
		data, err = json.Marshal(notes)
	}

	if err != nil {
		log.Println("I got an error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
	http.HandleFunc("/notes/", notesHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
