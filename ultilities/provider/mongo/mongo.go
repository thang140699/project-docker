package mongo

import (
	"fmt"
	"log"

	"WeddingBackEnd/ultilities/mongo"
)

type MongoProvider struct {
	mongoClient *mongo.MongoDB
}

func NewMongoProviderFromURL(u string) *MongoProvider {
	client := mongo.NewMongoDBFromURL(u)
	if client == nil {
		log.Fatalf("Mongo server connected unsuccessfully %s", u)
	}

	return &MongoProvider{
		mongoClient: client,
	}
}

func NewMongoProvider(server, user, password, database, collection, source string) *MongoProvider {
	configMongo := make(map[string]string)
	configMongo[mongo.DB_SERVER] = server
	configMongo[mongo.DB_USERNAME] = user
	configMongo[mongo.DB_PASSWORD] = password
	configMongo[mongo.DB_DATABASE] = database
	configMongo[mongo.DB_COLLECTION] = collection
	if user != "" {
		configMongo[mongo.DB_SOURCE] = source
	}

	client := mongo.NewMongoDB(configMongo)
	if client == nil {
		log.Fatalf("Mongo server connected unsuccessfully: nill client")
	}

	return &MongoProvider{
		mongoClient: client,
	}
}

func (provider *MongoProvider) MongoClient() *mongo.MongoDB {
	return provider.mongoClient
}

func (provider *MongoProvider) NewError(e error) error {
	if e == nil {
		return nil
	}
	return DatabaseExecutionError{
		Err:     e,
		Message: fmt.Sprintf("Mongo execution error: %s", e.Error()),
	}
}

type DatabaseExecutionError struct {
	Err     error
	Message string
}

func (e DatabaseExecutionError) Error() string {
	return e.Message
}

func (e DatabaseExecutionError) Unwrap() error {
	return e.Err
}
