package bookings

import "os"

// We define configuration default for running locally and overridable
// via enviroment variables so the application can be configured easily
// in the docker compose.

var (
	POSTGRES_HOST = GetOrElse("POSTGRES_HOST", "localhost")
	POSTGRES_PORT = GetOrElse("POSTGRES_PORT", "5431")

	PORT          = GetOrElse("PORT", "3001")
	APARTMENT_URL = GetOrElse("APARTMENT_URL", "http://localhost:3000")

	MQ_CONNECTION_STRING           = GetOrElse("MQ_CONNECTION_STRING", "amqp://guest:guest@localhost:5672/")
	MQ_APPARTMENT_CREATED_EXCHANGE = "appartment_created"
	MQ_APPARTMENT_CREATED_QUEUE    = "appartment_created"
	MQ_APPARTMENT_DELETED_EXCHANGE = "appartment_deleted"
	MQ_APPARTMENT_DELETED_QUEUE    = "appartment_deleted"
)

func GetOrElse(key string, d string) string {
	var value = os.Getenv(key)
	if value == "" {
		return d
	}
	return value
}