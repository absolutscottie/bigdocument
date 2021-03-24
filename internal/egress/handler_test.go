package egress

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/absolutscottie/bigdocument/internal/data"
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
			method:         "GET",
			expectedResult: false,
		},
		HandlersTestCase{
			name:           "no matching method",
			url:            "http://localhost:8181/document/test",
			method:         "POST",
			expectedResult: false,
		},
		HandlersTestCase{
			name:           "no matching path",
			url:            "http://localhost:8181/document/test",
			method:         "GET",
			expectedResult: true,
		},
	}

	router := mux.NewRouter()
	AddHandlers(router)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req, _ := http.NewRequest(testCase.method, testCase.url, nil)
			var match mux.RouteMatch
			result := router.Match(req, &match)
			if result != testCase.expectedResult {
				t.Errorf("Unexpected match result: %v\n", result)
			}
		})
	}
}

func TestHandleGetDocument(t *testing.T) {
	datastore, _ := data.NewMongoDatastore()
	ConfigureDatastore(datastore)

	router := mux.NewRouter()
	AddHandlers(router)

	datastore.DeleteDocument("test")
	testDoc, _ := datastore.NewDocument("test")
	testDoc.AddWord("the")
	testDoc.AddWord("quick")
	testDoc.AddWord("brown")
	testDoc.AddWord("fox")
	testDoc.AddWord("jumped")
	testDoc.AddWord("over")
	testDoc.AddWord("the")
	testDoc.AddWord("lazy")
	testDoc.AddWord("dog")

	req, err := http.NewRequest("GET", "http://localhost:8181/document/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	resultWords := make(map[string]bool)

	bodyBytes, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(bodyBytes))
	for scanner.Scan() {
		line := scanner.Text()
		if _, ok := resultWords[line]; ok {
			t.Errorf("duplicate word found in response: %v\n", line)
		}
		resultWords[line] = true
	}

	if len(resultWords) != 8 {
		t.Fatalf("Unexpected number of words returned: %d", len(resultWords))
	}
}
