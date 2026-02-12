// Notes API - Step 1
package main

import "fmt"

type Note struct {
	ID 		int
	Title 	string
	Content string
}

var store = map[string]Note{}


func main() {
	store["1"] = Note{ ID: 1, Title: "First Note", Content: "Just a note"}
	fmt.Println("Notes API")
}
