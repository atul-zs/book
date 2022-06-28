package Author

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

type Author struct {
	AuthorId  int    `json:"authorId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Dob       string `json:"dob"`
	PenName   string `json:"penName"`
}

type Book struct {
	BookId        string  `json:"bookId"`
	AuthorId      int     `json:"authorId"`
	Title         string  `json:"title"`
	Publication   string  `json:"publication"`
	PublishedDate string  `json:"publishedDate"`
	Author        *Author `json:"author"`
}

var db *sql.DB

func ConnectDb() *sql.DB {

	db, err := sql.Open("mysql", "root:7071686616@/Project1")
	if err != nil {
		log.Fatal("error in opening db")
	}
	if db.Ping() != nil {
		log.Fatal("connection to database failed!")
	}
	fmt.Println("connected")
	return db
}

func FetchAllBooks() []Book {
	Db := ConnectDb()
	defer Db.Close()

	Rows, err := Db.Query("SELECT * FROM book")

	if err != nil {
		fmt.Errorf("%v\n", err)
	}
	defer Rows.Close()

	var bk []Book

	for Rows.Next() {
		var b Book
		err := Rows.Scan(&b.BookId, &b.AuthorId, &b.Title, &b.Publication, &b.PublishedDate)
		if err != nil {
			fmt.Errorf("%v\n", err)
		}

		_, author := FetchAuthor(b.AuthorId)
		b.Author = &author
		bk = append(bk, b)
	}

	return bk
}
func FetchAuthor(id int) (int, Author) {
	Db := ConnectDb()
	defer Db.Close()

	Row := Db.QueryRow("SELECT * FROM author where author_id=?", id)
	var author Author
	if err := Row.Scan(&author.AuthorId, &author.FirstName, &author.LastName, &author.Dob, &author.PenName); err != nil {
		fmt.Errorf("failed: %v\n", err)
	}
	return author.AuthorId, author
}

func GetAllBooks(w http.ResponseWriter, req *http.Request) {

	books := FetchAllBooks()

	if req.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	mBook, err := json.Marshal(books)
	if err != nil {
		fmt.Errorf("%v\n", err)
	}
	bytes.NewBuffer(mBook)

	_, err = w.Write(mBook)
	if err != nil {
		log.Print(err)
	}

}

func GetBookById(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)

	if req.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	strings.ToLower(params["id"])
	id, _ := strconv.Atoi(params["id"])
	if id <= 0 {
		log.Print("invalid id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Db := ConnectDb()

	row := Db.QueryRow("select * from book where book_id=?", params["id"])
	var b Book
	if err := row.Scan(&b.BookId, &b.AuthorId, &b.Title, &b.Publication, &b.PublishedDate); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, author := FetchAuthor(b.AuthorId)
	b.Author = &author

	data, err := json.Marshal(b)
	if err != nil {
		log.Print(err)
		return
	}
	w.Write(data)
	w.WriteHeader(http.StatusOK)

}

func checkDob(Dob string) bool {

	dob := strings.Split(Dob, "/")
	day, _ := strconv.Atoi(dob[0])
	month, _ := strconv.Atoi(dob[1])
	year, _ := strconv.Atoi(dob[2])

	switch {
	case day <= 0 || day > 31:
		return false
	case month <= 0 || month > 12:
		return false
	case year > 2010:
		return false
	}
	return true
}

func PostAuthor(w http.ResponseWriter, req *http.Request) {
	body := req.Body

	data, err := io.ReadAll(body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
	}

	var author Author
	json.Unmarshal(data, &author)

	a, _ := FetchAuthor(author.AuthorId)
	if a == author.AuthorId || author.FirstName == "" {
		log.Print("error")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !checkDob(author.Dob) {
		log.Print("not valid Dob")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Db := ConnectDb()

	_, err = Db.Exec("insert into author(author_id,first_name,last_name,dob, pen_name)values(?,?,?,?,?)", author.AuthorId,
		author.FirstName, author.LastName, author.Dob, author.PenName)
	if err != nil {
		log.Print(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func checkPublishDate(PublishDate string) bool {
	p := strings.Split(PublishDate, "/")
	day, _ := strconv.Atoi(p[0])
	month, _ := strconv.Atoi(p[1])
	year, _ := strconv.Atoi(p[2])

	switch {
	case day < 0 || day > 31:
		return false
	case month < 0 || month > 12:
		return false
	case year > 2022 || year < 1880:
		return false
	}

	return true
}

func checkPublication(publication string) bool {
	strings.ToLower(publication)

	if !(publication == "penguin" || publication == "scholastic" || publication == "arihant") {
		return false
	}
	return true
}

func PostBook(w http.ResponseWriter, req *http.Request) {

	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	var book Book
	json.Unmarshal(body, &book)

	if id, _ := strconv.Atoi(book.BookId); id <= 0 {
		log.Print("invalid entry")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if book.BookId == "" || book.AuthorId <= 0 || book.Author.FirstName == "" || book.Title == "" {
		log.Print("not valid entry")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !checkDob(book.Author.Dob) || !checkPublishDate(book.PublishedDate) || !checkPublication(book.Publication) {
		log.Print("not valid entry")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Db := ConnectDb()

	res := Db.QueryRow("select * from book where book_id=?", book.BookId)
	var checkExitingId Book
	_ = res.Scan(&checkExitingId.BookId)
	if checkExitingId.BookId == book.BookId {
		log.Print("failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	a, _ := FetchAuthor(book.AuthorId)
	if a != book.AuthorId {
		log.Print("author does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = Db.Exec("insert into book(book_id,author_id,title,publication,publication_date)values (?,?,?,?,?)", book.BookId,
		book.AuthorId, book.Title, book.Publication, book.PublishedDate)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//w.Write(body)
	w.WriteHeader(http.StatusCreated)
}

func DeleteBook(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Errorf("invalid id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if id <= 0 {
		fmt.Println("invalid id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Db := ConnectDb()
	_ = Db.QueryRow("delete from book where book_id=?", params["id"])
	w.WriteHeader(http.StatusNoContent)
}

func DeleteAuthor(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Errorf("invalid id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if id <= 0 {
		fmt.Println("invalid id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Db := ConnectDb()
	_ = Db.QueryRow("delete from author where author_id=?", params["id"])
	w.WriteHeader(http.StatusNoContent)
}

func PutAuthor(w http.ResponseWriter, req *http.Request) {

	Data, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Errorf("failed:%v\n", err)
		return
	}
	var author Author
	json.Unmarshal(Data, &author)

	params := mux.Vars(req)
	Db := ConnectDb()

	if !checkDob(author.Dob) {
		fmt.Println("no valid Dob")
		w.WriteHeader(http.StatusBadRequest)
	}

	id, _ := strconv.Atoi(params["id"])

	var checkAuthor Author
	row := Db.QueryRow("select * from author where author_id=?", id)

	if err = row.Scan(&checkAuthor.AuthorId, &checkAuthor.FirstName, &checkAuthor.LastName, &checkAuthor.Dob, &checkAuthor.PenName); err == nil {
		fmt.Println(checkAuthor)
		_ = Db.QueryRow("delete from author where author_id=?", checkAuthor.AuthorId)
		_, err = Db.Exec("insert into author(author_id,first_name,last_name,dob, pen_name)values(?,?,?,?,?)",
			author.AuthorId, author.FirstName, author.LastName, author.Dob, author.PenName)

		w.WriteHeader(http.StatusCreated)
		w.Write(Data)
	} else {
		fmt.Println(checkAuthor)
		_, err = Db.Exec("insert into author(author_id,first_name,last_name,DOB, pen_name)values(?,?,?,?,?)", author.AuthorId, author.FirstName, author.LastName, author.Dob, author.PenName)

		w.WriteHeader(http.StatusCreated)
		w.Write(Data)
	}
}
