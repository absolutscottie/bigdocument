package data

type Document interface {
	AddWord(word string) error
	Read(b []byte) (int, error)
}

type Datastore interface {
	NewDocument(name string) (Document, error)
	FindDocument(name string) (Document, error)
	DeleteDocument(name string) error
}
