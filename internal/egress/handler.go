package egress

import (
	"net/http"

	"github.com/absolutscottie/bigdocument/internal/data"
	"github.com/gorilla/mux"
)

var (
	datastore data.Datastore
)

func ConfigureDatastore(ds data.Datastore) {
	datastore = ds
}

func AddHandlers(router *mux.Router) {
	router.HandleFunc("/document/{document_name}", handleGetDocument).Methods("GET")
}

func handleGetDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["document_name"]
	if name == "" {
		// there needs to be a name to find the document
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	doc, err := datastore.FindDocument(name)
	if err != nil {
		//could be a couple of different things but let's assume the file didn't exist
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for {
		line, err := doc.ReadLine()
		if err != nil {
			break
		}
		w.Write([]byte(line))
	}
}
