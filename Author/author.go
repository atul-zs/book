package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

type Book struct {
	BookId        int    `json:"bookId"`
	Title         string `json:"title"`
	Author        Author `json:"author"`
	Publication   string `json:"publication"`
	PublishedDate string `json:"publishedDate"`
}

type Author struct {
	AuthorId  int    `json:"authorId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Dob       string `json:"dob"`
	PenName   string `json:"penName"`
}

var db *sql.DB

func connectDb() *sql.DB {

	db, err := sql.Open("mysql", "root:7071686616@/Info")
	if err != nil {
		log.Fatal("error in opening db")
	}
	if db.Ping() != nil {
		log.Fatal("connection to database failed!")
	}
	fmt.Println("connected")
	return db
}
func main() {

	r := mux.NewRouter()
	r.HandleFunc("/book", GetAllBook)
	//r.HandleFunc("/book/", GetBookById)
	//r.HandleFunc("http://localhost:8000/book", PostBook)
	//r.HandleFunc("http://localhost:8000/author", PostAuthor)

	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal(err)
	}

}

func PostAuthor(w http.ResponseWriter, req *http.Request) {

	db = connectDb()
	var author *Author
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
		return
	}

	err = json.Unmarshal(body, &author)
	if err != nil {
		log.Print(err)
		return
	}

	if author.FirstName == "" || author.LastName == "" || author.PenName == "" {
		log.Print("mention the name")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(author.Dob) < 10 || len(author.Dob) > 10 {
		log.Print("invalid format")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := db.Query("SELECT author_id from author WHERE first_name=? AND last_name=? AND dob=? AND pen_name=?", author.FirstName, author.LastName, author.Dob, author.PenName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !res.Next() || err != nil {
		_, err = db.Exec("INSERT INTO author (first_name,last_name,dob,pen_name) VALUES (?,?,?,?)", author.FirstName, author.LastName, author.Dob, author.PenName)
		if err != nil {
			log.Print(err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
	w.Write(body)
	w.WriteHeader(http.StatusCreated)

}

func FetchAuthor(id int) (int, Author) {
	db := connectDb()
	defer db.Close()

	row := db.QueryRow("select * from author where author_id=?", id)

	var author Author
	err := row.Scan(&author.AuthorId, &author.FirstName, &author.LastName, &author.Dob, &author.PenName)
	if err != nil {
		log.Print(err)
	}
	return author.AuthorId, author

}

func GetAllBook(w http.ResponseWriter, req *http.Request) {

	db := connectDb()
	defer db.Close()

	rows, err := db.Query("select * from book")
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()

	var book []Book
	for rows.Next() {
		var b Book
		err := rows.Scan(&b.BookId, &b.Title, &b.Author, &b.Publication, &b.PublishedDate)
		if err != nil {
			log.Print(err)
		}
		_, author := FetchAuthor(b.Author.AuthorId)
		//fmt.Println(author)
		b.Author = author
		book = append(book, b)
	}

	data, err := json.Marshal(book)
	if err != nil {
		log.Print(err)
	}
	bytes.NewBuffer(data)

	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}

//func GetBookById(w http.ResponseWriter, req *http.Request) {
//
//	db := connectDb()
//	defer db.Close()
//
//	rows, err := db.Query("select * from book")
//	if err != nil {
//		log.Print(err)
//	}
//	defer rows.Close()
//
//	var book []Book
//	for rows.Next() {
//		var b Book
//		err := rows.Scan(&b.Id, &b.Title, &b.Author, &b.Publication, &b.PublishedDate)
//		if err != nil {
//			log.Print(err)
//		}
//		_, author := FetchAuthor(b.Author.AuthorId)
//		//fmt.Println(author)
//		b.Author = author
//		book = append(book, b)
//	}
//
//	params := mux.Vars(req)
//
//	for _, item := range book {
//
//		if item.Id == params["id"] {
//
//			data, err := json.Marshal(item)
//			if err != nil {
//				log.Print(err)
//			}
//
//			bytes.NewBuffer(data)
//			_, err = w.Write(data)
//			if err != nil {
//				log.Print(err)
//			}
//		}
//	}
//
//}

func CheckPublishedDate(s string) bool {
	publicationDate := strings.Split(s, "/")
	if len(publicationDate) < 3 {
		log.Print("invalid  publication date ")
		return false
	}
	year, _ := strconv.Atoi(publicationDate[2])

	if year < 1800 || year > 2022 {
		log.Print("invalid Publication ")
		return false
	}
	return true
}
func PostBook(w http.ResponseWriter, req *http.Request) {

	var book Book
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &book)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if book.Publication != "Scholastic" && book.Publication != "Penguin" && book.Publication != "Arihant" {
		log.Print("invalid entry")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !CheckPublishedDate(book.PublishedDate) {
		log.Print("invalid publication date")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if book.Title == "" || book.Author.FirstName == "" {
		log.Print("invalid entry")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if book.BookId <= 0 {
		log.Print("invalid bookid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	db = connectDb()
	var id int64
	row, err := db.Query("SELECT author_id from author WHERE first_name=? AND last_name=? AND dob=? AND pen_name=?", book.Author.FirstName, book.Author.LastName, book.Author.Dob, book.Author.PenName)
	if row.Next() {
		row.Scan(&id)
	}
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO book (book_id,title,author_id,publication,published_date) VALUES (?,?,?,?,?)", book.BookId, book.Title, id, book.Publication, book.PublishedDate)
	fmt.Println()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write(body)
	w.WriteHeader(http.StatusOK)
}
