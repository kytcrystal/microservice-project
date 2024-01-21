function getOrElse(key, defaultValue) {
    return process.env[key] ? process.env[key] : defaultValue;
}

const PORT          = getOrElse("PORT", "3002")
const APARTMENT_URL = getOrElse("APARTMENT_URL", "http://localhost:3000") 
const BOOKING_URL = getOrElse("BOOKING_URL", "http://localhost:3001") 
const MQ_CONNECTION_STRING          = getOrElse("MQ_CONNECTION_STRING", "amqp://guest:guest@localhost:5672/")
const MQ_APARTMENT_CREATED_EXCHANGE = "apartment_created"
const MQ_APARTMENT_CREATED_QUEUE    = "search-service.apartment_created"
const MQ_APARTMENT_DELETED_EXCHANGE = "apartment_deleted"
const MQ_APARTMENT_DELETED_QUEUE    = "search-service.apartment_deleted"
const MQ_BOOKING_CREATED_EXCHANGE   = "booking_created"
const MQ_BOOKING_CREATED_QUEUE      = "search-service.booking_created"
const MQ_BOOKING_CANCELLED_EXCHANGE = "booking_cancelled"
const MQ_BOOKING_CANCELLED_QUEUE    = "search-service.booking_cancelled"
const MQ_BOOKING_UPDATED_EXCHANGE   = "booking_updated"
const MQ_BOOKING_UPDATED_QUEUE      = "search-service.booking_updated"

module.exports = {
    PORT,
    APARTMENT_URL,
    BOOKING_URL,
    MQ_CONNECTION_STRING,
    MQ_APARTMENT_CREATED_EXCHANGE,
    MQ_APARTMENT_CREATED_QUEUE,
    MQ_APARTMENT_DELETED_EXCHANGE,
    MQ_APARTMENT_DELETED_QUEUE,
    MQ_BOOKING_CREATED_EXCHANGE,
    MQ_BOOKING_CREATED_QUEUE,
    MQ_BOOKING_CANCELLED_EXCHANGE,
    MQ_BOOKING_CANCELLED_QUEUE,
    MQ_BOOKING_UPDATED_EXCHANGE,
    MQ_BOOKING_UPDATED_QUEUE
};
