package Author

import (
	"bytes"
	_ "bytes"
	"encoding/json"
	"fmt"
	_ "fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"

	_ "github.com/gorilla/mux"
)

func TestGetAllBooks(t *testing.T) {

	Testcases := []struct {
		desc        string
		methodInput string
		target      string
		expOut      []Book
		status      int
	}{
		{"success case:", "GET", "/book", []Book{{"10",
			5, "story", "scholastic", "29/08/2010", &Author{5, "Kyrie",
				"Erving", "19/05/1994", "KE"}},
			{"14", 7, "story", "scholastic", "29/08/2010", &Author{7, "Jasum",
				"tatum", "19/04/1996", "JT"}},
		}, http.StatusOK,
		},
		{"success case:", "POST", "/book", []Book{{"10",
			5, "story", "scholastic", "29/08/2010", &Author{5, "Kyrie",
				"Erving", "19/05/1994", "KE"}},
			{"14", 7, "story", "scholastic", "29/08/2010", &Author{7, "Jasum",
				"tatum", "19/04/1996", "JT"}},
		}, http.StatusBadRequest,
		},
	}

	for _, tc := range Testcases {
		req := httptest.NewRequest(tc.methodInput, "http://localhost:8000"+tc.target, nil)
		w := httptest.NewRecorder()
		GetAllBooks(w, req)

		resp := w.Result()
		body, _ := io.ReadAll(resp.Body)

		var book []Book

		_ = json.Unmarshal(body, &book)

		if resp.StatusCode != tc.status {
			t.Errorf("failed for %s", tc.desc)

			if reflect.DeepEqual(tc.expOut, book) {
				t.Errorf("%v", book)
			}

		}

	}
}

func TestGetBookById(t *testing.T) {
	Testcases := []struct {
		desc               string
		methodInput        string
		bookId             string
		expected           Book
		expectedStatusCode int
	}{
		{"success case:", "GET", "11",
			Book{"11", 5, " story", "scholastic", "29/08/2010",
				&Author{5, "Kyrie", "Erving", "19/05/1994", "KE"}}, http.StatusOK},
		{"invalid method", "POST", "2",
			Book{}, http.StatusBadRequest},
		{"invalid id (negative)", "GET", "-2",
			Book{"-2", 1, "the cup", "penguin", "10/07/2019",
				&Author{1, "atul", "gond", "29/06/2000", "ag"}}, http.StatusBadRequest},
	}

	for _, tc := range Testcases {

		req := httptest.NewRequest(tc.methodInput, "http://localhost:8000/book/{id}"+tc.bookId, nil)
		w := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"id": tc.bookId})
		GetBookById(w, req)

		resp := w.Result()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var book Book
		_ = json.Unmarshal(body, &book)

		if resp.StatusCode != tc.expectedStatusCode {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}

func TestPostAuthor(t *testing.T) {
	testcases := []struct {
		desc        string
		inputMethod string
		target      string
		body        Author
		expected    int
	}{
		{"successful case:", "POST", "/author", Author{
			20, "John", "Dere", "19/04/1981", "JD"}, http.StatusCreated},

		{" author already exists", "POST", "/author", Author{
			3, "Kevin", "Durant", "19/05/1995", "KD"}, http.StatusBadRequest},

		{"invalid firstname", "POST", "/author", Author{
			3, "", "mrinal", "20/05/1990", "KD"}, http.StatusBadRequest},

		{"invalid Dob", "POST", "/author", Author{
			3, "nilotpal", "mrinal", "96/23/2000", "KD"}, http.StatusBadRequest},
	}

	for _, tc := range testcases {

		author, err := json.Marshal(tc.body)
		if err != nil {
			fmt.Println("error:", err)
		}

		req := httptest.NewRequest(tc.inputMethod, "https://localhost:8000"+tc.target, bytes.NewBuffer(author))
		w := httptest.NewRecorder()
		PostAuthor(w, req)

		res := w.Result()
		if res.StatusCode != tc.expected {
			t.Errorf("failed for %s", tc.desc)
		}
	}
}

func TestPostBook(t *testing.T) {

	testcases := []struct {
		desc        string
		inputMethod string
		target      string
		body        Book
		expected    int
	}{
		{"valid case", "POST", "/book", Book{"15", 7, "go",
			"scholastic", "29/08/2002", &Author{7, "Jasum",
				"Tatum", "19/04/1996", "JT"}},
			http.StatusOK},

		{"invalid author DOb", "POST", "/book", Book{"4", 23, "story",
			"penguin", "20/03/2020", &Author{23, "ram",
				"kumar", "00/22/2020", "rk"}},
			http.StatusBadRequest},

		{"invalid bookId", "POST", "/book", Book{"-4", 14, "story",
			"penguin", "23/02/2012", &Author{14, "atul",
				"kumar", "29/06/2000", "Ak"}},
			http.StatusBadRequest},

		{"invalid author's firstName", "POST", "/book", Book{"4", 8, "story",
			"scholastic", "20/03/2010", &Author{8, "",
				"Erving", "19/05/1994", "KE"}},
			http.StatusBadRequest},

		{"not existing author", "POST", "/book", Book{"5", 1, "story",
			"arihant", "20/03/2010", &Author{1, "atul",
				"kumar", "30/00/2001", "ak"}},
			http.StatusBadRequest},

		{"invalid publication", "POST", "/book", Book{"6", 3, "story",
			"sun", "20/03/2010", &Author{3, "Atul",
				"kumar", "30/00/2002", "Ak"}},
			http.StatusBadRequest},

		{"invalid title", "POST", "/book", Book{"7", 5, "",
			"penguin", "20/03/2010", &Author{5, "shani",
				"kumar", "30/00/2001", "sk"}},
			http.StatusBadRequest},

		{"invalid publishedDate", "POST", "/book", Book{"8", 6, "story",
			"McGrowHill", "20/03/1789", &Author{6, "Atul",
				"kumar", "30/00/2001", "Ak"}},
			http.StatusBadRequest},
	}

	for _, tc := range testcases {

		b, err := json.Marshal(tc.body)
		if err != nil {
			fmt.Println("error:", err)
		}

		req := httptest.NewRequest(tc.inputMethod, "http://localhost:8000"+tc.target, bytes.NewBuffer(b))
		w := httptest.NewRecorder()
		PostBook(w, req)

		res := w.Result()
		fmt.Println(res.StatusCode)
		if res.StatusCode != tc.expected {
			t.Errorf("failed for %v\n", tc.desc)
		}
	}
}

func TestDeleteBook(t *testing.T) {
	testcases := []struct {
		desc        string
		inputMethod string
		target      string
		expected    int
	}{
		{"valid id", "DELETE", "15", http.StatusNoContent},
		{"invalid id", "DELETE", "-4", http.StatusBadRequest},
	}

	for _, tc := range testcases {

		req := httptest.NewRequest(tc.inputMethod, "https://localhost:8000/book/{id}"+tc.target, nil)
		w := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"id": tc.target})
		DeleteBook(w, req)

		res := w.Result()
		if res.StatusCode != tc.expected {
			t.Errorf("failed for %s", tc.desc)
		}

	}
}

func TestDeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc        string
		inputMethod string
		target      string
		expected    int
	}{
		{"valid authorId", "DELETE", "4", http.StatusNoContent},
		{"invalid authorId", "DELETE", "-3", http.StatusBadRequest},
	}

	for _, tc := range testcases {

		req := httptest.NewRequest(tc.inputMethod, "https://localhost:8000/author/{id}"+tc.target, nil)
		w := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{"id": tc.target})
		DeleteAuthor(w, req)

		res := w.Result()
		if res.StatusCode != tc.expected {
			t.Errorf("failed for %s", tc.desc)
		}
		if reflect.DeepEqual(tc.expected, res) {
			t.Errorf("%v", res)
		}
	}
}

func TestPutBook(t *testing.T) {
	testcases := []struct {
		desc        string
		inputMethod string
		target      string
		body        Book
		expected    int
	}{
		{"successful case", "PUT", "localhost", Book{"17", 9, "new",
			"arihant", "29/08/2010", &Author{9, "Salman",
				"Khan", "19/04/1980", "SK"}}, http.StatusCreated},
		{"invalid case: Get method instead of put", "PUT", "localhost", Book{"17", 9, "new",
			"arihant", "29/08/2010", &Author{9, "Salman",
				"Khan", "19/04/1980", "SK"}}, http.StatusBadRequest},
	}

	for _, tc := range testcases {

		//req := httptest.NewRequest(tc.inputMethod, tc.target, nil)
		w := httptest.NewRecorder()
		//PutBook(w, req)

		res := w.Result()
		if res.StatusCode != tc.expected {
			t.Errorf("failed for %s", tc.desc)
		}
		if reflect.DeepEqual(tc.expected, res.StatusCode) {
			t.Errorf("%v", res.StatusCode)
		}
	}
}

func TestPutAuthor(t *testing.T) {
	testcases := []struct {
		desc        string
		inputMethod string
		target      string
		author      Author
		expected    int
	}{
		{"successful case", "PUT", "local",
			Author{20, "ram", "kumar", "4/05/2004", "rk"}, http.StatusCreated},
		{"invalid case: put method is not mentioned", "GET", "local",
			Author{25, "ram", "kumar", "4/05/2004", "rk"}, http.StatusBadRequest},
		{"invalid case: ", "PUT", "local",
			Author{26, "sujit", "kumar", "4/05/2004", "sk"}, http.StatusBadRequest},
	}

	for _, tc := range testcases {

		req := httptest.NewRequest(tc.inputMethod, tc.target, nil)
		w := httptest.NewRecorder()
		PutAuthor(w, req)

		res := w.Result()
		if res.StatusCode != tc.expected {
			t.Errorf("failed for %s", tc.desc)
		}
		if reflect.DeepEqual(tc.expected, res.StatusCode) {
			t.Errorf(" failed %v", res.StatusCode)
		}
	}
}

// PutBook : updates the particular field in book table and if not exits then creates
func PutBook(w http.ResponseWriter, req *http.Request) {

	body := req.Body
	params := mux.Vars(req)
	ReqBody, err := io.ReadAll(body)
	if err != nil {
		log.Print(err)
		return
	}

	var book Book
	json.Unmarshal(ReqBody, &book)

	id, author := FetchAuthor(book.AuthorId)
	if id != book.AuthorId {
		log.Print("author does not exist")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	book.Author = &author

	Db := ConnectDb()

	if !checkPublishDate(book.PublishedDate) || !checkPublication(book.Publication) || book.Title == "" || !checkDob(book.Author.Dob) {
		log.Print("invalid entry")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var checkBook Book
	row := Db.QueryRow("select * from book where book_id=?", params["id"])

	if err = row.Scan(&checkBook.BookId, &checkBook.AuthorId, &checkBook.Title, &checkBook.Publication, &checkBook.PublishedDate); err == nil {

		_, err = Db.Exec("update book set bookId=?,authorId=?,title=?,publication=?,publishedDate=? where bookId=?",
			book.BookId, book.AuthorId, book.Title, book.Publication, book.PublishedDate, params["id"])

		fmt.Fprintf(w, "successfull updated id =%v\n", params["id"])
		w.WriteHeader(http.StatusCreated)
		return

	} else {
		_, err = Db.Exec("insert into book(book_id,author_id,title,publication,published_date)values(?,?,?,?,?) ",
			book.BookId, book.AuthorId, book.Title, book.Publication, book.PublishedDate)

		fmt.Fprintf(w, "successfull inserted id =%v\n", params["id"])
		w.WriteHeader(http.StatusCreated)
		return
	}

}
