# RabbitMQ

RabbitMQ is the open source message broker that implement amqp protocol

It has 3 components

1. The RabbitMQ server: The message broker itself. It's a server you can run locally (or on a server in production) to send and receive messages.
2. Client library: We'll be programming in Go, so our code will use the AMQP library to interact with the RabbitMQ server.
3. Management UI: RabbitMQ comes with a management UI that you can use to monitor and manage your RabbitMQ server. It's a web-based UI that you can access in your browser.

Run this command in a terminal to start RabbitMQ:

```sh
docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.13-management
```

This will download and run the RabbitMQ server and ManagementUI until you kill it with Ctrl+C
The server itself is running on port 5672, but you won't interact with it in your browser, just through code. However, you can use the management UI in your browser. Open http://localhost:15672 and log in with the username guest and password guest.

# AMQP Protocol

- A communication protocol is just some rules that define how messages are sent, received, and understood.
- Just like we define how can is, then so we can communicate with the term cat to convey the four leg monster meaning

# MQTT and STOMP

- MQTT is designed for small IoT devices and as such is optimized to be lightweight and energy efficient.
- STOMP is designed for web applications and is designed to be simple and easy to use.
- AMQP is a sophisticated, feature-rich protocol that's great for complex routing, filtering, and delivery requirements.

Assignment:

Configure our RabbitMQ instance to support STOMP as well, just for fun. We'll need to create our own Dockerfile that installs the RabbitMQ STOMP plugin.

1. Create a file called Dockerfile with the following contents:

```docker
FROM rabbitmq:3.13-management
RUN rabbitmq-plugins enable rabbitmq_stomp
```

2, Build the image, and name it rabbitmq-stomp:

```sh
docker build -t rabbitmq-stomp .
```

3. If you have rabbit running, use ./rabbit.sh stop to stop it.
4. Take a look at the contents of rabbit.sh... you should see the line that's used to run a new RabbitMQ container with docker run. Copy that line into your terminal, but change it to do the following:
   1. Use the image you just built
   2. Expose the STOMP port: 61613
   3. Run the container.

Once it's running, test port 61613 with NetCat to make sure that it's open and that RabbitMQ is running the STOMP plugin:

```sh
nc -vz localhost 61613
```

You should get a "connection succeeded" message if all is well.

```sh
docker run -d --name rabbitmq_stomp -p 61613:61613 -p 15672:15672 rabbitmq-stomp
docker build -t rabbitmq-stomp .
```

# Exchanges and Queues

In RabbitMQ, an exchange is where publishers send messages, typically with a routing key.
The exchange takes the message, uses the routing key as a filter, and sends the message to any queues that are listening for that routing key.
Publishers don't know about queues at all. They just send messages to exchanges, sometimes with a routing key.
![alt text](image.png)

- Exchange: A routing agent that sends messages to queues.
- Binding: A link between an exchange and a queue that uses a routing key to decide which messages go to the queue.
- Queue: A buffer in the RabbitMQ server that holds messages until they are consumed.
- Channel: A virtual connection inside a connection that allows you to create queues, exchanges, and publish messages.
- Connection: A TCP connection to the RabbitMQ server.

Assingment:
Assignment
Let's update our server to publish pause/resume messages to an exchange on a specific routing key. The server can then communicate with all the various players of the game to let them know when the game is paused or resumed. We'll handle the queues and consumption of the messages later.

1. In cmd/server/main.go, after opening the RabbitMQ connection, create a new channel using the .Channel method on the connection.
2. Create a new package: internal/pubsub. This is where we'll put reusable code for interacting with RabbitMQ, that way we can use it in both the server and client.
3. Create an exported PublishJSON function in the internal/pubsub package. Here's its signature:

```go
func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error
```

4. In PublishJSON :
   1. Marshal the val to JSON bytes
   2. Use the channel's .PublishWithContext method to publish the message to the exchange with the routing key. Some configurations:
      1. Set ctx to context.Background()
      2. Set mandatory to false.
      3. Set immediate to false.
      4. In the amqp.Publishing struct, you only need to set two fields:
         1. ContentType to "application/json".
         2. Body to the JSON bytes.

5. Back in the server code, use the PublishJSON function to publish a message to the exchange!
   1. Use the channel you created.
   2. Use the internal/routing package's ExchangePerilDirect string for the exchange.
   3. Use the internal/routing package's PauseKey string for the routing key.
   4. Pass a PlayingState value from the internal/routing package with IsPaused set to true.
6. Make sure the RabbitMQ Docker container is running in the background, if it's not run ./rabbit.sh start.
7. Run your server with:

```sh
go run ./cmd/server
```

8. Watch the logs (./rabbit.sh logs) of the RabbitMQ server... you should see an error like this:
   no exchange 'peril_direct' in vhost '/'
9. Go to the RabbitMQ management UI at http://localhost:15672 and navigate to the "Exchanges" tab. Create a new exchange called peril_direct with the type direct.
10. Rerun the server. You should see the message get published without any errors in the RabbitMQ logs.
    While there are no hard errors, the message will be "unroutable" because there are no queues bound to the exchange yet, but we'll fix that later.

# Types of Exchanges

RabbitMQ supports several types of exchanges, each serving a different routing strategy.
![alt text](image-1.png)
In my experience, direct and topic are the most commonly useful in backend Pub/Sub architectures. I rarely have a use for sending all messages to all queues, or for routing based on the message headers.

Breakdown:

1. Direct: Messages are routed to the queues based on the message routing key exactly matching the binding key of the queue.
2. Topic: Messages are routed to queues based on wildcard matches between the routing key and the routing pattern specified in the binding.
3. Fanout: It routes messages to all of the queues bound to it, ignoring the routing key.
4. Headers: Routes based on header values instead of the routing key. It's similar to topic but uses message header attributes for routing.

# Create a Queue

- Queues are where the messages are stored after being routed through the exchange. Messages sit in a queue until they are consumed by a subscriber.

Durability

- Queues can be "durable" or "transient". Durable queues survive a RabbitMQ server restart, while transient queues do not.
- The metadata of a durable queue is stored on disk, while transient queues are only stored in memory.

Assignment:
Let's manually create a queue to capture the "pause" messages our server is sending.

1. Open the Management UI and navigate to the "Queues and Streams" tab. Click the "Add a new queue" dropdown at the bottom left.
2. Name the queue pause_test, because we're just going to use it temporarily to test our server's ability to publish messages.
3. Leave the durability as "Durable" and create the queue. If you screw it up, you can click on the queue, delete it, and try again.
4. Click on the queue, then go to the "Bindings" section. You'll be able to see that the queue is already bound to the default exchange. Add another binding to the peril_direct exchange.
   1. "From exchange": peril_direct
   2. "Routing key": Use the exact string of the PauseKey constant from the internal/routing package. This is a direct exchange, so the routing key must match exactly.
5. Click "Bind".
6. Restart your server to publish a new message to the exchange.
7. If everything worked, you should see a new "ready" message populate in the queue in the management UI. It's just sitting there patiently waiting to be consumed.
8. After selecting the queue, scroll down to the "Get messages" tab. Click "Get Message(s)" to see the message that was published by the server. The payload is the JSON representation of the PlayingState struct with the IsPaused field set to true, though the UI may display it as an encoded string.
9. The message should still be in the queue because "Nack message requeue true" will put the message back after showing it to you.

# Transient Queues

Let's update our code to automatically create and delete transient queues, rather than doing it manually. We'll create the queues that the client will use to receive the "pause" messages from the server.

## Durable and Transient Queue Types

Durable queues survive a RabbitMQ server restart, while transient queues do not. We can also set the auto-delete and exclusive properties of our queues:

- Exclusive: The queue can only be used by the connection that created it.
- Auto-delete: The queue will be automatically deleted when its last connection is closed.

For simplicity of our game, we'll make our transient and durable queues always have the same properties:

- "Durable" queues in our system will always be non-exclusive and non-auto-delete.
- "Transient" queues will always be exclusive and auto-delete.

## Assignment

1. Update the `cmd/client` package to connect to Rabbit, similar to the `cmd/server` package.
2. Use the `ClientWelcome()` function in `internal/gamelogic` to prompt the user for a username.
3. Declare and bind a transient queue by creating and using a new function in the `internal/pubsub` package. I called mine `DeclareAndBind` with this signature:

```go
func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error)
```

4. In `DeclareAndBind`:
   1. Create a new `.Channel()` on the connection.
   2. Declare a new queue using `.QueueDeclare()`:
      1. The `durable` parameter should only be true if `queueType` is durable.
      2. The `autoDelete` parameter should be true if `queueType` is transient.
      3. The `exclusive` parameter should be true if `queueType` is transient.
      4. The `noWait` parameter should be false.
      5. The `args` parameter should be nil.

   3. Bind the queue to the exchange using `.QueueBind()`.
   4. Return the channel and queue.

5. Back in the `cmd/client` package, use these parameters to call `DeclareAndBind`:
   1. `exchange`: `peril_direct` (this is a constant in the `internal/routing` package)
   2. `queueName`: `pause.username` where `username` is the user's input. The `pause` section of the name is the routing key constant in the `internal/routing` package and is joined by a `.`.
   3. `routingKey`: `pause` (this is a constant in the `internal/routing` package)
   4. `queueType`: transient

6. After declaring and binding the queue, the client should wait for a Ctrl+C signal to exit.
7. Run the client! Enter `suntzu` (all lowercase) as the username.
8. It will print out some nonsense about possible commands, but you can ignore that. You can't run them yet because we haven't implemented them.
9. Check the RabbitMQ management UI to see if the queue was created and bound to the exchange.
10. Close the client. If all goes well, the queue should automatically be deleted after a few seconds.
11. Run the client as `suntzu` again, making sure the queue is recreated.

# Decoupling

- One of the big advantages of a Pub/Sub architecture over a point-to-point messaging system is decoupling.
- Our server publishes the "Hey, I'm pausing the game" message _once_ to the broker, and anyone who cares can get a copy of the message in their own queue, all without the server knowing or caring who's listening.

## Assignment

First, let's update the server so that it can interactively pause and resume the game.

1. Run the `PrintServerHelp` function in `internal/gamelogic` as the server starts up so that you can see the commands the user of the REPL can use.
2. Start an infinite loop.
3. At the beginning of the loop, use the `GetInput` function in `internal/gamelogic` to wait for a slice of input "words" from the user. If the slice is empty, continue to the next iteration of the loop.
4. Check the first word:
   1. If it's `"pause"`, log to the console that you're sending a pause message, and publish the pause message as you were doing before.
   2. If it's `"resume"`, log to the console that you're sending a resume message, and publish the resume message as you were doing before. The only difference is that the `IsPaused` field should be set to `false`.
   3. If it's `"quit"`, log to the console that you're exiting, and break out of the loop.
   4. If it's anything else, log to the console that you don't understand the command.
5. Test the app:
   1. Start 3 clients, each in their own terminal. Use the usernames `washington`, `napoleon`, and `churchill`.
   2. Run the server again to pause the game.
   3. You should see 3 queues created in the RabbitMQ management UI, each with its own copy of the "pause" message.
