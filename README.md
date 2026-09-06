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
