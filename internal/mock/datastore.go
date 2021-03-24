package mock

import (
	"errors"
	"fmt"
	"io"

	"github.com/absolutscottie/bigdocument/internal/data"
)

type MockDocument struct {
	words        map[string]bool
	internalBuf  []byte
	internalUsed int
	needNewline  bool
}

func (m *MockDocument) AddWord(word string) error {
	m.words[word] = true
	return nil
}

func (m *MockDocument) ReadLine() (string, error) {
	if len(m.words) == 0 {
		return "", io.EOF
	}

	for k, _ := range m.words {
		return fmt.Sprintf("%s\n", k), nil
	}

	return "", io.EOF
}

func (m *MockDocument) Count() (int64, error) {
	return int64(len(m.words)), nil
}

type MockDatastore struct {
	documents map[string]*MockDocument
}

func NewDatastore() *MockDatastore {
	return &MockDatastore{
		documents: make(map[string]*MockDocument),
	}
}

func (d *MockDatastore) NewDocument(name string) (data.Document, error) {
	document := &MockDocument{
		words: make(map[string]bool),
	}

	d.documents[name] = document
	return document, nil
}

func (d *MockDatastore) FindDocument(name string) (data.Document, error) {
	if _, ok := d.documents[name]; !ok {
		return nil, errors.New("document not found")
	} else {
		return d.documents[name], nil
	}
}

func (d *MockDatastore) DeleteDocument(name string) error {
	delete(d.documents, name)
	return nil
}
