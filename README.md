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
