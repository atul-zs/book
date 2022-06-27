package book

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

//const Publications string= "Scholastic" " Arihant" "Penguin"
//
//func TestGetAllBook(t *testing.T) {
//	testCases := []struct {
//		desc      string
//		input     string
//		expOut    []Book
//		expStatus int
//	}{
//		{"success case:", "book",
//			[]Book{{1, "the wall", Author{10, "shani", "kumar", "10/02/2001", "Mark Twain"}, "Scolastic", "23/10/2018"},
//				{2, "the beast", Author{11, "sujit", "kumar", "12/02/2001", "Lewis Carroll"}, "Arihant", "01/10/2021"},
//				{3, "the evolution", Author{12, "shiv", "kumar", "22/02/2002", "Sapphire"}, "Penguin", "02/02/2017"}}, http.StatusOK},
//		{"failure case:", "book",
//			[]Book{}, http.StatusBadRequest},
//	}
//
//	for _, tc := range testCases {
//		res := httptest.NewRequest("GET", "http://localhost:8000/"+tc.input, nil)
//		w := httptest.NewRecorder()
//		GetAllBook(w, res)
//		r := w.Result()
//
//		data, err := ioutil.ReadAll(res.Body)
//		if err != nil {
//			log.Print(err)
//		}
//
//		var book []Book
//		err = json.Unmarshal(data, &book)
//		if err != nil {
//			t.Errorf("expected error to be nil got %v", err)
//		}
//
//		if r.StatusCode != tc.expStatus {
//			t.Errorf("invalid :%v", r.StatusCode)
//		}
//		// assert.Equal(tc.expOut,book)
//		if reflect.DeepEqual(tc.expOut, book) {
//			t.Errorf(" got %v", book)
//		}
//	}
//
//}

//func TestGetBookId(t *testing.T) {
//	testCases := []struct {
//		desc      string
//		input     string
//		expOut    Book
//		expStatus int
//	}{
//		{"Success case:",
//			"2",
//			Book{1, "the beast", Author{10, "shani", "kumar", "10/02/2001", "Mark Twain"}, "Arihant", "01/10/2021"}, 200,
//		},
//		{"Invalid input:",
//			"a",
//			Book{}, 404,
//		},
//		{"Id does not exist in database :",
//			"100",
//			Book{}, 404,
//		},
//	}
//
//	for _, tc := range testCases {
//		res := httptest.NewRequest("GET", "http://localhost:8000/book/"+tc.input, nil)
//		w := httptest.NewRecorder()
//		GetBookId(w, res)
//		r := w.Result()
//		defer r.Body.Close()
//
//		data, err := ioutil.ReadAll(res.Body)
//		if err != nil {
//			t.Errorf("expected error to be nil got %v", err)
//		}
//		var book Book
//		err = json.Unmarshal(data, &book)
//		if err != nil {
//			t.Errorf("expected error to be nil got %v", err)
//		}
//
//		if r.StatusCode != 200 {
//			t.Errorf("invalid :%s", r.StatusCode)
//		}
//		if reflect.DeepEqual(tc.expOut, book) {
//			t.Errorf(" got %v", book)
//		}
//	}
//
//}
//
func TestPostBook(t *testing.T) {
	testCases := []struct {
		desc      string
		input     string
		inputBody Book
		expOut    int
	}{
		{"1: Publication Arihant", "/book",
			Book{51, "the river", Author{FirstName: "john", LastName: "carter", Dob: "20/12/2010", PenName: "Lemony Snicket"}, "Arihant", "12/06/2017"},
			http.StatusOK,
		},
		{"2: Publication Scholastic ", "/book",
			Book{32, "the state", Author{FirstName: "dhiran", LastName: "kumar", Dob: "12/03/2000", PenName: "sun"}, "Scholastic", "12/06/1900"},
			http.StatusOK,
		},
		{" 3: Missing title ", "/book",
			Book{34, "", Author{FirstName: "dhiran", LastName: "kumar", Dob: "12/03/2000", PenName: "sun"}, "Scholastic", "12/06/1900"},
			http.StatusBadRequest,
		},
		{" 4:Invalid published date", "/book",
			Book{35, "the lion", Author{FirstName: "sujit", LastName: "kumar", Dob: "12/03/2000", PenName: "flower"}, "Penguin", "12/06/1730"},
			http.StatusBadRequest,
		},
		{"5:Invalid published date", "/book",
			Book{36, "the elephant", Author{FirstName: "sanket", LastName: "kumar", Dob: "12/03/2000", PenName: "moon"}, "Penguin", "12/06/2023"},
			http.StatusBadRequest,
		},
		{"6:Missing first name of author ", "/book",
			Book{37, "the knight", Author{FirstName: "", LastName: "kumar", Dob: "12/03/2000", PenName: "cactus"}, "Penguin", "12/08/2000"},
			http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {

		data, err := json.Marshal(tc.inputBody)
		if err != nil {
			t.Errorf("ERROR: %v", err)
		}

		res, err := http.NewRequest("POST", "http://localhost:8000"+tc.input, bytes.NewBuffer(data))
		if err != nil {
			t.Errorf("ERROR:bt %v", err)
		}
		w := httptest.NewRecorder()
		PostBook(w, res)
		r := w.Result()

		if r.StatusCode != tc.expOut {
			t.Errorf("Error %v ", r.StatusCode)
		}

	}

}

//func TestPostAuthor(t *testing.T) {
//	testCases := []struct {
//		desc      string
//		input     string
//		inputBody Author
//		expOut    int
//	}{
//		{"test 1: ", "/author",
//			Author{FirstName: "Kevin", LastName: "Durant", Dob: "20/12/2013", PenName: "shark"},
//			http.StatusCreated,
//		},
//		{"test 2: Author last name not entered", "/author",
//			Author{FirstName: "kevin", LastName: "", Dob: "29/03/2000", PenName: "12/06/2017"},
//			http.StatusBadRequest,
//		},
//		{"test 3: Successful entry", "/author",
//			Author{FirstName: "dhiran", LastName: "kumar", Dob: "12/03/2000", PenName: "sun"},
//			http.StatusCreated,
//		},
//		{"test 4:Invalid  dob", "/author",
//			Author{FirstName: "krie", LastName: "erving", Dob: "22/34", PenName: "whale"},
//			http.StatusBadRequest,
//		},
//		{"test 5: PenName not entered", "/author",
//			Author{FirstName: "kris", LastName: "paul", Dob: "12/02/2012", PenName: ""},
//			http.StatusBadRequest,
//		},
//		{"test 5: Successful entry", "/author",
//			Author{FirstName: "jordan", LastName: "krish", Dob: "12/02/2000", PenName: "thor"},
//			http.StatusBadRequest,
//		},
//	}
//
//	for _, tc := range testCases {
//
//		input, err := json.Marshal(tc.inputBody)
//
//		if err != nil {
//			t.Errorf("ERROR: %v", err)
//		}
//
//		res, err := http.NewRequest("POST", "http://localhost:8000"+tc.input, bytes.NewBuffer(input))
//		if err != nil {
//			t.Errorf("Error")
//		}
//		w := httptest.NewRecorder()
//		PostAuthor(w, res)
//		r := w.Result()
//
//		if r.StatusCode != tc.expOut {
//			t.Errorf("status not found")
//		}
//
//	}
//
//}

//func TestPut(t *testing.T){
//	testCases:=[]struct {
//
//	}
//
//}
//func TestDelete(t *testing.T) {
//
//}
