const Database = require("better-sqlite3");
const db = new Database("search.db", { verbose: console.log });
const express = require("express");
const amqp = require("amqplib");

const tables = require("./tables");
const config = require("./config");

const app = express();

app.get("/", (req, res) => {
  res.send("Hello, World!");
});

app.get("/api/search/available", (req, res) => {
  const from = req.query.from || "";
  const to = req.query.to || "";
  const row = tables.searchAvailableApartments(db, from, to);
  res.send(row);
});

app.get("/api/search/apartments", (req, res) => {
  const row = tables.listAll(db, "apartments");
  res.send(row);
});

app.get("/api/search/bookings", (req, res) => {
  const row = tables.listAll(db, "bookings");
  res.send(row);
});

async function startListener() {
	console.log("Starting Listener For Apartment Messages")

  const connection = await amqp.connect(config.MQ_CONNECTION_STRING);
  const channel = await connection.createChannel();

  await messageReceiver(channel, config.MQ_APARTMENT_CREATED_EXCHANGE, 
    config.MQ_APARTMENT_CREATED_QUEUE, tables.createApartment);

  await messageReceiver(channel, config.MQ_APARTMENT_DELETED_EXCHANGE, 
    config.MQ_APARTMENT_DELETED_QUEUE, tables.deleteApartment);

  await messageReceiver(channel, config.MQ_BOOKING_CREATED_EXCHANGE, 
    config.MQ_BOOKING_CREATED_QUEUE, tables.createBooking);

  await messageReceiver(channel, config.MQ_BOOKING_CANCELLED_EXCHANGE, 
    config.MQ_BOOKING_CANCELLED_QUEUE, tables.cancelBooking);

  await messageReceiver(channel, config.MQ_BOOKING_UPDATED_EXCHANGE, 
    config.MQ_BOOKING_UPDATED_QUEUE, tables.updateBooking);
}

tables.createTable(db);
startListener();

// Start the server and listen for incoming requests
app.listen(config.PORT, () => {
  console.log(`Server listening on http://localhost:${config.PORT}`);
});

async function messageReceiver(channel, exchange, queue, actOnMessage) {
  
  await channel.assertExchange(exchange, 'fanout', {
    durable: true
  });

  await channel.assertQueue(queue, {
    exclusive: true
  })

  await channel.bindQueue(queue, exchange, '');

  channel.consume(queue, function(msg) {
    try {
    if(msg.content) {
        console.log("Received a message: '%s'", msg.content.toString());
        actOnMessage(db, JSON.parse(msg.content.toString()))
      }
    } catch(err) {
      console.error(`error while processing message from ${queue}: ${err}`, msg)
    }
  }, {
    noAck: true
  });
}
