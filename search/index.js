const Database = require("better-sqlite3");
const db = new Database("search.db", { verbose: console.log });
const express = require("express");
const amqp = require("amqplib");

const tables = require("./tables");

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

	const RABBIT_MQ_CONNECTION_STRING = "amqp://guest:guest@localhost:5672/"
  const connection = await amqp.connect(RABBIT_MQ_CONNECTION_STRING);
  const channel = await connection.createChannel();

  const MQ_APARTMENT_CREATED_EXCHANGE = "apartment_created";
	const MQ_APARTMENT_CREATED_QUEUE    = "search-service.apartment_created";
  await messageReceiver(channel, MQ_APARTMENT_CREATED_EXCHANGE, MQ_APARTMENT_CREATED_QUEUE, tables.createApartment);

  const MQ_APARTMENT_DELETED_EXCHANGE = "apartment_deleted";
	const MQ_APARTMENT_DELETED_QUEUE    = "search-service.apartment_deleted";
  await messageReceiver(channel, MQ_APARTMENT_DELETED_EXCHANGE, MQ_APARTMENT_DELETED_QUEUE, tables.deleteApartment);

  const MQ_BOOKING_CREATED_EXCHANGE = "booking_created";
	const MQ_BOOKING_CREATED_QUEUE    = "search-service.booking_created";
  await messageReceiver(channel, MQ_BOOKING_CREATED_EXCHANGE, MQ_BOOKING_CREATED_QUEUE, tables.createBooking);

  const MQ_BOOKING_CANCELLED_EXCHANGE = "booking_cancelled";
	const MQ_BOOKING_CANCELLED_QUEUE    = "search-service.booking_cancelled";
  await messageReceiver(channel, MQ_BOOKING_CANCELLED_EXCHANGE, MQ_BOOKING_CANCELLED_QUEUE, tables.cancelBooking);

  const MQ_BOOKING_UPDATED_EXCHANGE = "booking_updated";
	const MQ_BOOKING_UPDATED_QUEUE    = "search-service.booking_updated";
  await messageReceiver(channel, MQ_BOOKING_UPDATED_EXCHANGE, MQ_BOOKING_UPDATED_QUEUE, tables.updateBooking);
}

tables.createTable(db);
startListener();

const port = 3002;

// Start the server and listen for incoming requests
app.listen(port, () => {
  console.log(`Server listening on http://localhost:${port}`);
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
    if(msg.content) {
        console.log("Received a message: '%s'", msg.content.toString());
        actOnMessage(db, JSON.parse(msg.content.toString()))
      }
  }, {
    noAck: true
  });
}
