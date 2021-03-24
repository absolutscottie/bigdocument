package ingest

import (
	"bufio"
	"io"
	"net/http"

	"github.com/absolutscottie/bigdocument/internal/data"
	"github.com/gorilla/mux"
)

var (
	datastore data.Datastore
)

// ConfigureDatastore set's the package-level datastore to the provided datastore
func ConfigureDatastore(ds data.Datastore) {
	datastore = ds
}

// AddHandler configures the provided router with the needed methods to support
// ingest operations
func AddHandlers(router *mux.Router) {
	router.HandleFunc("/document/{document_name}", handlePutDocument).Methods("PUT")
}

func handlePutDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	defer r.Body.Close()

	docName := vars["document_name"]
	if docName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Delete any document that previously existed with this same name
	err := datastore.DeleteDocument(docName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create a new document
	document, err := datastore.NewDocument(docName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = readWords(r.Body, document)
	if err != nil {
		// The assumption here is that something was wrong with the input
		// probably not a good assumption.
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//200 OK is written by default
}

// readWords reads output from the provided Reader, separated by new lines. Each
// arbitrary length 'word' read from the Reader is stored in the provided document
func readWords(body io.Reader, document data.Document) error {
	// scanner was chosen here to easily identify line breaks in the input stream
	// it seemed like the best choice without giving it too much thought.
	// It has the added benefit of stripping off the newlines
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		word := scanner.Text()
		err := document.AddWord(word)
		if err != nil {
			return err
		}
	}
	return scanner.Err()
}
