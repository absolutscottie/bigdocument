package mock

import (
	"errors"

	"github.com/absolutscottie/bigdocument/internal/common/data"
)

type MockDocument struct {
	words map[string]bool
}

func (m *MockDocument) AddWord(word string) error {
	m.words[word] = true
	return nil
}

func (m *MockDocument) Read(b []byte) (int, error) {
	return 0, nil
}

type MockDatastore struct {
	documents map[string]*MockDocument
}

func NewDatastore() MockDatastore {
	return MockDatastore{
		documents: make(map[string]*MockDocument),
	}
}

func (d MockDatastore) NewDocument(name string) (data.Document, error) {
	document := &MockDocument{
		words: make(map[string]bool),
	}

	d.documents[name] = document
	return document, nil
}

func (d MockDatastore) FindDocument(name string) (data.Document, error) {
	if _, ok := d.documents[name]; !ok {
		return nil, errors.New("document not found")
	} else {
		return d.documents[name], nil
	}
}

func (d MockDatastore) DeleteDocument(name string) error {
	if _, ok := d.documents[name]; !ok {
		return errors.New("document not found")
	}
	delete(d.documents, name)
	return nil
}
