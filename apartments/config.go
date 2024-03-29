package apartments

import "os"

// We define configuration default for running locally and overridable
// via enviroment variables so the application can be configured easily
// in the docker compose.

var (
	POSTGRES_HOST = GetOrElse("POSTGRES_HOST", "localhost")
	POSTGRES_PORT = GetOrElse("POSTGRES_PORT", "5432")

	PORT = GetOrElse("PORT", "3000")

	MQ_CONNECTION_STRING          = GetOrElse("MQ_CONNECTION_STRING", "amqp://guest:guest@localhost:5672/")
	MQ_APARTMENT_CREATED_EXCHANGE = "apartment_created"
	MQ_APARTMENT_DELETED_EXCHANGE = "apartment_deleted"
)

func GetOrElse(key string, d string) string {
	var value = os.Getenv(key)
	if value == "" {
		return d
	}
	return value
}
