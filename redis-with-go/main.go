package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"github.com/go-redis/redis/v8"
)

func main() {
    // Print a message to the console.
    fmt.Println("Connecting to Redis...")

    // Connect to Redis.
    client := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", 		// Use container name as hostname.
        Password: "",                         // No password set.
        DB:       0,                          // Use the default DB.
    })

    // Ping the Redis server to check if the connection is established.
    pong, err := client.Ping(context.Background()).Result()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Connected to Redis:", pong)

	// Define a struct.
	type Person struct {
		Name string `json:"name"`
		Age int		`json:"age"`
		Occupation string `json:"occupation"`
	}
	
	// Set a key-value pair with a struct.
	jsonString, err := json.Marshal(Person{
		"xuan", 
		28, 
		"developer",
	})
	if err != nil {
		log.Fatal(err)
	}

    // // Set a key-value pair.
    // err = client.Set(context.Background(), "name", "xuan", 0).Err()
    // if err != nil {
    //     log.Fatal(err)
	// }
	// fmt.Println("Set key-value pair:", "name", "xuan")

    // // Get the value of a key.
    // val, err := client.Get(context.Background(), "name").Result()
    // if err != nil {
    //     log.Fatal(err)
	// }
	// fmt.Printf("Get value of key:, %s\n", val)

    // Set a key-value pair with a struct.
    err = client.Set(context.Background(), "person", jsonString, 0).Err()
    if err != nil {
        log.Fatal(err)
	}
	// fmt.Println("Set key-value pair:", "person", jsonString)

    // Get the value of a key.
    val, err := client.Get(context.Background(), "person").Result()
    if err != nil {
        log.Fatal(err)
	}
	fmt.Printf("Get value of key:, %s\n", val)

	// Unmarshal the JSON string into a Person struct.
	var person Person
	if err := json.Unmarshal([]byte(val), &person); err != nil {
		log.Fatal(err)
	}

	// Access the "occupation" field.
	occupation := person.Occupation
	fmt.Println("Occupation:", occupation)

    // Close the connection when you're done.
    defer client.Close()
}