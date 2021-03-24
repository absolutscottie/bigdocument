package mock

import (
	"errors"
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

func (m *MockDocument) Read(b []byte) (int, error) {
	//return 0, nil
	used := 0
	for k, _ := range m.words {
		if m.internalBuf == nil {
			delete(m.words, k)
			m.internalBuf = []byte(k)
		}

		for m.internalUsed < len(m.internalBuf) && used < len(b) {
			if m.needNewline {
				m.needNewline = false
				b[used] = '\n'
				used++
				continue
			}

			b[used] = m.internalBuf[m.internalUsed]
			used++
			m.internalUsed++
		}

		if m.internalUsed >= len(m.internalBuf) {
			m.internalBuf = nil
			m.internalUsed = 0
		}

		if used >= len(b) {
			return used, nil
		}

		m.needNewline = true
	}
	return used, io.EOF
}

func (m *MockDocument) Count() int {
	return len(m.words)
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
