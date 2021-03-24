package ingest

import (
	"net/http"
	"strings"
	"testing"

	"github.com/absolutscottie/bigdocument/internal/mock"
	"github.com/gorilla/mux"
)

type HandlersTestCase struct {
	name           string
	url            string
	method         string
	expectedResult bool
}

func TestAddHandlers(t *testing.T) {
	testCases := []HandlersTestCase{
		HandlersTestCase{
			name:           "no matching path",
			url:            "http://localhost:8181/documents/test",
			method:         "PUT",
			expectedResult: false,
		},
		HandlersTestCase{
			name:           "no matching method",
			url:            "http://localhost:8181/document/test",
			method:         "GET",
			expectedResult: false,
		},
		HandlersTestCase{
			name:           "no matching path",
			url:            "http://localhost:8181/document/test",
			method:         "PUT",
			expectedResult: true,
		},
	}

	router := mux.NewRouter()
	AddHandlers(router)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req, _ := http.NewRequest(testCase.method, testCase.url, strings.NewReader("onetwothree"))
			var match mux.RouteMatch
			result := router.Match(req, &match)
			if result != testCase.expectedResult {
				t.Errorf("Unexpected match result: %v\n", result)
			}
		})
	}
}

type ReadWordsTestCase struct {
	name          string
	content       string
	expectedCount int
}

func TestReadWords(t *testing.T) {
	testCases := []ReadWordsTestCase{
		ReadWordsTestCase{
			name:          "test 1",
			expectedCount: 8,
			content: `the
quick
brown
fox
jumped
over
the
lazy
dog`,
		},
		ReadWordsTestCase{
			name:          "test 2",
			expectedCount: 1,
			content: `word
word
word
word
word`,
		},
		ReadWordsTestCase{
			name:          "test 3",
			expectedCount: 0,
			content:       "",
		},
	}
	for _, tc := range testCases {
		ds := mock.NewDatastore()
		t.Run(tc.name, func(t *testing.T) {
			doc, _ := ds.NewDocument(tc.name)
			err := readWords(strings.NewReader(tc.content), doc)
			if err != nil {
				t.Fatalf("Unexpected error: %v\n", err)
			}

			if doc.Count() != tc.expectedCount {
				t.Fatalf("Unexpected number of words counted - expected %d but found %d\n", tc.expectedCount, doc.Count())
			}
		})
	}
}
