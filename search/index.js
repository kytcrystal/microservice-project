const express = require('express');
const app = express();

app.get('/', (req, res) => {
  res.send('Hello, World!');
});

app.get('/api/search/available', (req, res) => {
    const from = req.query.from || '';
    const to = req.query.to || '';
  
    // Send a response with the parsed query parameter
    res.send(`Hello, searching for date from ${from} to ${to}!`);
  });

// Specify the port to listen on
const port = 3002;

// Start the server and listen for incoming requests
app.listen(port, () => {
  console.log(`Server listening on http://localhost:${port}`);
});
