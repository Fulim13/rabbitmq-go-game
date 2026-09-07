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

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
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

	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalln("failed to create channel")
	}

	_, _, err = pubsub.DeclareAndBind(conn, "peril_topic", "game_logs", "game_logs.*", pubsub.Durable)
	if err != nil {
		log.Fatalln("failed to create channel")
	}

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Publishing paused game state")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: true,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "resume":
			fmt.Println("Publishing resumes game state")
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{
					IsPaused: false,
				},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		case "quit":
			log.Println("goodbye")
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
