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

  process.once("SIGINT", async () => {
    await channel.close();
    await connection.close();
  });

	const APARTMENTS_CREATED_QUEUE = "apartment_created"
  await messageReceiver(channel, APARTMENTS_CREATED_QUEUE, tables.createApartment);
	const APARTMENTS_DELETED_QUEUE = "apartment_deleted"
  await messageReceiver(channel, APARTMENTS_DELETED_QUEUE, tables.deleteApartment);
}



tables.createTable(db);
startListener();

const port = 3002;

// Start the server and listen for incoming requests
app.listen(port, () => {
  console.log(`Server listening on http://localhost:${port}`);
});

async function messageReceiver(channel, queue, actOnMessage) {
  await channel.assertQueue(queue, { durable: false });
  await channel.consume(
    queue,
    (message) => {
      if (message) {
        console.log(
          "Received a message: '%s'",
          JSON.parse(message.content.toString())
        );
        actOnMessage(db, JSON.parse(message.content.toString()))
      }
    },
    { noAck: true }
  );
}
