package main

import (
	"Golearn/book/Author"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/book", Author.GetAllBooks).Methods("GET")
	r.HandleFunc("/book/{id}", Author.GetBookById).Methods("GET")
	r.HandleFunc("/book", Author.PostBook).Methods(http.MethodPost)
	r.HandleFunc("/author", Author.PostAuthor).Methods(http.MethodPost)
	r.HandleFunc("/book/{id}", Author.DeleteBook).Methods("DELETE")
	r.HandleFunc("/author/{id}", Author.DeleteAuthor).Methods("DELETE")
	r.HandleFunc("/author/{id}", Author.PutAuthor).Methods("PUT")
	//r.HandleFunc("/author/{id}", Author.PutBook).Methods("PUT")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected")

}
