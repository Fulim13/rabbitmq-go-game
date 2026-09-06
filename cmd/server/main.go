package main

// The internal/gamelogic package which contains the game logic for the Peril game.
// The internal/routing package which contains some routing constants for the game.
// Stubs of the cmd/client and cmd/server packages, which are main packages that run the client and server for the game.
// The rabbit.sh script: a convenience for starting and stopping RabbitMQ with Docker. You can run:
// ./rabbit.sh start to start RabbitMQ
// ./rabbit.sh stop to stop it
// ./rabbit.sh logs to view the server logs

// Practice 1:
// Install the Go AMQP library:
// go get github.com/rabbitmq/amqp091-go

// Update the cmd/server package's main.go. You'll want to use an alias to import the amqp package:
// import (
//     amqp "github.com/rabbitmq/amqp091-go"
// )

// In the main function, do the following:
// 1. Declare a connection string with the value: amqp://guest:guest@localhost:5672/. This is how your application will know where to connect to the RabbitMQ server.
// 2. Call amqp.Dial with the connection string to create a new connection to RabbitMQ.
// 3. Defer a .Close() of the connection to ensure it's closed when the program exits.
// 4. Print a message to the console that the connection was successful.
// 5. Wait for a signal (e.g. Ctrl+C) to exit the program.
// 6. If a signal is received, print a message to the console that the program is shutting down and close the connection.

// To test your code, use the script to start a RabbitMQ server with Docker in the background:
// ./rabbit.sh start

// Once it's running, open a new terminal and run your server:
// go run ./cmd/server
// You should see your message that the connection was successful. Press Ctrl+C to exit the program and see the shutdown message.

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
