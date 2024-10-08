package databases

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
}

func Connect() (*MongoDB, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017" // Значение по умолчанию
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Проверка подключения
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return &MongoDB{Client: client}, nil
}

func (database *MongoDB) Disconnect() error {
	return database.Client.Disconnect(context.TODO())
}

func (database *MongoDB) Collection(name string) *mongo.Collection {
	return database.Client.Database("testdb").Collection(name)
}
