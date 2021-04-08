package data

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/absolutscottie/bigdocument/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DatabaseName = "words"
)

// MongoDocument probably a bad name, but MongoDocument is a wrapper for a mongo collection
// that implements the Document inteface
type MongoDocument struct {
	name       string
	collection *mongo.Collection

	skip int
}

// AddWord will add the provided string to the wrapped collection
func (m *MongoDocument) AddWord(word string) error {
	_, err := m.collection.InsertOne(context.TODO(), bson.D{
		{Key: "text", Value: word},
	})
	return err
}

// Count will return the number of documents in the wrapped collection
func (m *MongoDocument) Count() (int64, error) {
	filter := bson.M{"document": m.name}
	result, err := m.collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return -1, err
	}
	return result, nil
}

func (m *MongoDocument) ReadLine() (string, error) {
	found, err := m.findWords()
	if err != nil {
		return "", err
	}
	if len(found) == 0 {
		return "", io.EOF
	}
	return fmt.Sprintf("%s\n", found[0]), nil
}

func (m *MongoDocument) findWords() ([]string, error) {
	type Word struct {
		Text string `bson:"text"`
	}

	found := make([]string, 0)
	options := options.Find()
	options.SetLimit(1)
	options.SetSkip(int64(m.skip))
	m.skip++

	cursor, err := m.collection.Find(context.TODO(), bson.D{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var word Word
		if err = cursor.Decode(&word); err != nil {
			return nil, err
		}
		found = append(found, word.Text)
	}
	return found, nil
}

//	MongoDatastore is essentiall a wrapper for a mongodb client which implements the Datastore interface
type MongoDatastore struct {
	client *mongo.Client
}

//	NewMongoDatastore establishes a connection to a local mongodb install and returns a
//	wrapper for the client
func NewMongoDatastore(cfg *config.Config) (Datastore, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DatastoreHost))
	if err != nil {
		return nil, err
	}

	return MongoDatastore{
		client: client,
	}, nil
}

// NewDocument creates a new collection with the name provided. An index is added on 'text' and
// is considered unique to prevent duplicate entries
func (m MongoDatastore) NewDocument(name string) (Document, error) {
	mod := mongo.IndexModel{
		Keys: bson.M{
			"text": 1,
		},
		Options: options.Index().SetUnique(true),
	}
	collection := m.client.Database(DatabaseName).Collection(name)
	_, err := collection.Indexes().CreateOne(context.TODO(), mod)
	if err != nil {
		return nil, err
	}
	doc := &MongoDocument{
		name:       name,
		collection: collection,
	}
	return doc, nil
}

// FindDocument determines whether there were any documents associated with the named
// collection. if none are found we consider that to be 'not found'
func (m MongoDatastore) FindDocument(name string) (Document, error) {
	collection := m.client.Database(DatabaseName).Collection(name)
	return &MongoDocument{
		name:       name,
		collection: collection,
	}, nil
}

// DeleteDocument drops the collection with the provided name
func (m MongoDatastore) DeleteDocument(name string) error {
	err := m.client.Database(DatabaseName).Collection(name).Drop(context.TODO())
	return err
}
