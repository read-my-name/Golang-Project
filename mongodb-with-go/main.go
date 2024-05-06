package main

import (
    "context"
    "fmt"
    "log"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
    // Set client options.
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

    // Connect to MongoDB.
    client, err := mongo.Connect(context.Background(), clientOptions)
    if err != nil {
        log.Fatal(err)
    }

    // Check the connection.
    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to MongoDB!")

	// Create a reference to the "users" collection within the "testing" database.
    usersCollection := client.Database("testing").Collection("users")

    // Insert a document into the "users" collection.
    _, err = usersCollection.InsertOne(context.Background(), bson.M{"name": "Alice", "age": 30})
    if err != nil {
        log.Fatal(err)
    }

    // Query the database to retrieve the inserted document.
    var result bson.M
    err = usersCollection.FindOne(context.Background(), bson.M{"name": "Alice"}).Decode(&result)
    if err != nil {
        log.Fatal(err)
    }

    // Print the retrieved document.
    fmt.Println("Inserted Document:", result)

    // Disconnect from MongoDB.
    err = client.Disconnect(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Disconnected from MongoDB!")
}