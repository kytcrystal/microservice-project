package search

// We can use this file for some constants that could be duplicated in the two apps, so it's easier to copy and paste
// as well it's more difficult to do typos.

const APARTMENTS_QUEUE_NAME = "apartments-queue"

// TODO: make this configurable via env variable so it's easy to Dockerize the app
const RABBIT_MQ_CONNECTION_STRING = "amqp://guest:guest@localhost:5672/"
