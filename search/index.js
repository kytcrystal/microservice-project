const Database = require("better-sqlite3");
const db = new Database("search.db", { verbose: console.log });
const express = require("express");
const tables = require("./tables");

const app = express();

app.get("/", (req, res) => {
  res.send("Hello, World!");
});

app.get("/api/search/available", (req, res) => {
  const from = req.query.from || "";
  const to = req.query.to || "";

  res.send(`Hello, searching for date from ${from} to ${to}!`);
});

app.get("/api/search/apartments", (req, res) => {
  var row = tables.listAll(db, "apartments");
  res.send(row);
});

app.get("/api/search/bookings", (req, res) => {
  var row = tables.listAll(db, "bookings");
  res.send(row);
});



tables.createTable(db);

const port = 3002;

// Start the server and listen for incoming requests
app.listen(port, () => {
  console.log(`Server listening on http://localhost:${port}`);
});
