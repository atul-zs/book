package book

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func testGetAllBook(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		expOut []Book
	}{
		{"test 1:", "http://localhost:8000/book",
			[]Book{{2, "the wall", "tommy", "the moon", "23/10/2018"},
				{3, "the beast", "sam", "the sun", "01/10/2022"},
				{4, "the book", "shani", "the moon", "02/02/2017"}}},
	}

	for _, tc := range testCases {
		res := httptest.NewRequest("GET", tc.input, nil)
		w := httptest.NewRecorder()
		GetAllBook(w, res)
		r := w.Result()
		defer r.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		var book []Book
		err = json.Unmarshal(data, &book)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		if reflect.DeepEqual(tc.expOut, book) {
			t.Errorf(" got %v", book)
		}
	}

}

func testGetBookId(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		expOut Book
	}{
		{"test 1:",
			"http://localhost:8000/book/{5}",
			Book{5, "the new evolution", "tom", "the golden", "12/08/2014"},
		},
	}

	for _, tc := range testCases {
		res := httptest.NewRequest("GET", tc.input, nil)
		w := httptest.NewRecorder()
		GetBookId(w, res)
		r := w.Result()
		defer r.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}
		var book Book
		err = json.Unmarshal(data, &book)
		if err != nil {
			t.Errorf("expected error to be nil got %v", err)
		}

		if reflect.DeepEqual(tc.expOut, book) {
			t.Errorf(" got %v", book)
		}
	}

}

func testPostBook(t *testing.T) {
	testCases := []struct {
		desc      string
		input     string
		inputBody *Book
		expOut    int
	}{
		{"test 1: Publication Arihant", "http://localhost:8000/book/{21}",
			&Book{21, "the wall", "tommy", "Arihant", "12/06/2017"},
			200,
		},
		{"test 2: Invalid Publication", "http://localhost:8000/book/{22}",
			&Book{22, "the wall", "tommy", "the sun", "12/06/2017"},
			0,
		},
		{"test 3:Invalid published date", "http://localhost:8000/book/{23}",
			&Book{23, "the river", "john", "Penguin", "12/06/1879"},
			0,
		},
		{"test 4: Publication Penguin ", "http://localhost:8000/book/{24}",
			&Book{24, "the river", "john", "Penguin", "12/06/1900"},
			0,
		},
		{"test 5:Invalid published date", "http://localhost:8000/book/{25}",
			&Book{25, "the river", "john", "Penguin", "12/06/2023"},
			0,
		},
		{"test 6:Invalid Id", "http://localhost:8000/book/{100}",
			&Book{100, "little things", "john", "Penguin", "12/06/2020"},
			0,
		},
		{"test 7: publication name Scholastic", "http://localhost:8000/book/{27}",
			&Book{27, "monster", "pascal", "Scholastic", "12/06/2020"},
			200,
		},
	}

	for _, tc := range testCases {

		json_data, err := json.Marshal(tc.inputBody)

		if err != nil {
			t.Errorf("ERROR: %v", err)
		}

		res, err := http.NewRequest("POST", tc.input, bytes.NewBuffer(json_data))
		if err != nil {
			t.Errorf("ERROR:bt %v", err)
		}
		w := httptest.NewRecorder()
		PostBook(w, res)
		r := w.Result()

		if err != nil {
			t.Errorf("ERROR:bt %v", err)
		}
		fmt.Println(r)

	}

}

func testPostBookByAuthor(t *testing.T) {
	testCases := []struct {
		desc      string
		input     string
		inputBody *Book
		expOut    int
	}{
		{"test 2: Author id exist", "http://localhost:8000/author",
			&Book{21, "the wall", "tommy", "Arihant", "12/06/2017"},
			0,
		},
		{"test 2: Author Id does not exist", "http://localhost:8000/author",
			&Book{109, "the wall", "tommy", "the sun", "12/06/2017"},
			0,
		},
		{"test 3:Invalid published date", "http://localhost:8000/author",
			&Book{21, "the river", "john", "Penguin", "12/06/1879"},
			0,
		},
		{"test 3:Invalid published date", "http://localhost:8000/author",
			&Book{21, "the river", "john", "Penguin", "12/06/2023"},
			0,
		},
		{"test 7: publication name Scholastic", "http://localhost:8000/author",
			&Book{27, "bond", "pascal", "Scholastic", "12/06/2020"},
			200,
		},
	}

	for _, tc := range testCases {

		json_data, err := json.Marshal(tc.inputBody)

		if err != nil {
			t.Errorf("ERROR: %v", err)
		}

		res, err := http.NewRequest("POST", tc.input, bytes.NewBuffer(json_data))
		if err != nil {
			t.Errorf("ERROR:bt %v", err)
		}
		w := httptest.NewRecorder()
		PostBookByAuthorName(w, res)
		r := w.Result()

		if err != nil {
			t.Errorf("ERROR:bt %v", err)
		}
		fmt.Println(r)

	}

}
