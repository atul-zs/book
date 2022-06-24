package book

import "net/http"

type Book struct {
	Id            int
	Title         string
	Author        string
	Publication   string
	PublishedDate string
}

func GetBookId(w http.ResponseWriter, req *http.Request) {
	return
}

func GetAllBook(w http.ResponseWriter, req *http.Request) []Book {

	return []Book{}
}
func PostBook(w http.ResponseWriter, req *http.Request) {

}
func PostBookByAuthorName(w http.ResponseWriter, req *http.Request) {

}
