// Notes API - Step 6
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Note struct {
	ID      int
	Title   string
	Content string
}

var idCounter = 0

type ErrorResponse struct {
	Error string `json:"error"`
}

var store = map[int]Note{}

func notesHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		id := strings.Trim(strings.TrimPrefix(r.URL.Path, "/notes"), "/")
		if len(id) > 0 {
			getNoteById(id, w)
		} else {
			getNotes(w)
		}
	case http.MethodPost:
		note, err := readBody(r, w)
		if err != nil {
			return
		}
		addNote(note, w)
	case http.MethodPut:
		note, err := readBody(r, w)
		if err != nil {
			return
		}
		updateNote(note, w, r)
	case http.MethodDelete:
		deleteNote(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

}

func readBody(r *http.Request, w http.ResponseWriter) (Note, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("I got an error with POST request, reading body")
		w.WriteHeader(http.StatusBadRequest)
		return Note{}, err
	}
	var bodyRequest Note
	decodeError := json.Unmarshal(body, &bodyRequest)
	if decodeError != nil {
		log.Println("I got an error decoding body")
		w.WriteHeader(http.StatusBadRequest)
		return Note{}, decodeError
	}
	return bodyRequest, nil
}

func addNote(note Note, w http.ResponseWriter) {
	if note.ID == 0 {
		idCounter++
		note.ID = idCounter
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "ID not permitted"})
		return
	}
	store[idCounter] = note
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(note)
	w.Write(data)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	id := strings.Trim(strings.TrimPrefix(r.URL.Path, "/notes"), "/")
	value, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid ID"})
		return
	}
	if _, exists := store[value]; !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "not found"})
		return
	}
	delete(store, value)
	w.WriteHeader(http.StatusNoContent)
}

func getNotes(w http.ResponseWriter) {
	var data []byte
	var err error
	var notes []Note
	for _, note := range store {
		notes = append(notes, note)
	}
	data, err = json.Marshal(notes)
	if err != nil {
		log.Println("I got an error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func getNoteById(id string, w http.ResponseWriter) {
	var data []byte
	var err error
	value, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid ID"})
		return
	}
	if _, exists := store[value]; !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "not found"})
		return
	}
	data, err = json.Marshal(store[value])
	w.Write(data)
}

func updateNote(note Note, w http.ResponseWriter, r *http.Request) {
	id := strings.Trim(strings.TrimPrefix(r.URL.Path, "/notes"), "/")
	value, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid ID"})
		return
	}
	if _, exists := store[value]; !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "not found"})
		return
	}
	note.ID = value
	store[note.ID] = note
	data, _ := json.Marshal(store[note.ID])
	w.Write(data)
}

func main() {
	store[1] = Note{ID: 1, Title: "First Note", Content: "Just a note"}
	store[2] = Note{ID: 2, Title: "Second Note", Content: "It's getting interesting"}
	store[3] = Note{ID: 3, Title: "Third Note", Content: "Muy interesante!"}
	store[4] = Note{ID: 4, Title: "Fourth Note", Content: "Hola Amigo! Como estas?"}
	idCounter = 4
	fmt.Println("Notes API")

	fmt.Println("Server is running on port 8080")

	http.HandleFunc("/notes", notesHandler)
	http.HandleFunc("/notes/", notesHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
