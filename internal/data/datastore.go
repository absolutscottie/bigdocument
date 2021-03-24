package data

type Document interface {
	AddWord(word string) error
	Count() (int64, error)
	ReadLine() (string, error)
}

type Datastore interface {
	NewDocument(name string) (Document, error)
	FindDocument(name string) (Document, error)
	DeleteDocument(name string) error
}
