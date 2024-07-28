package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books []Book
var nextID int
var mutex sync.Mutex



func main() {
	
		router := mux.NewRouter()
// Initialize with some sample data
	books = append(books, Book{ID: 1, Title: "The Go Programming Language", Author: "Alan A. A. Donovan"})
	books = append(books, Book{ID: 2, Title: "Go in Action", Author: "William Kennedy"})
	nextID = 3
	
	
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books", createBook).Methods("POST")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")
	
		fmt.Println("Server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
	
}


func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for _, book := range books {
		if book.ID == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.NotFound(w, r)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	
	mutex.Lock()
	book.ID = nextID
	nextID++
	books = append(books, book)
	mutex.Unlock()

	json.NewEncoder(w).Encode(book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	var updatedBook Book
	_ = json.NewDecoder(r.Body).Decode(&updatedBook)

	mutex.Lock()
	for i, book := range books {
		if book.ID == id {
			updatedBook.ID = id
			books[i] = updatedBook
			json.NewEncoder(w).Encode(updatedBook)
			mutex.Unlock()
			return
		}
	}
	mutex.Unlock()
	http.NotFound(w, r)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	mutex.Lock()
	for i, book := range books {
		if book.ID == id {
			books = append(books[:i], books[i+1:]...)
			mutex.Unlock()
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	mutex.Unlock()
	http.NotFound(w, r)
}