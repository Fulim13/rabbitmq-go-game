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
