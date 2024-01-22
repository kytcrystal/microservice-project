package eventsourcing

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/jmoiron/sqlx"
)

type EventRepository struct {
	db *sqlx.DB
}

var POSTGRES_HOST = GetOrElse("POSTGRES_HOST", "localhost")
var POSTGRES_PORT = GetOrElse("POSTGRES_PORT", "5431")

func GetOrElse(key string, d string) string {
	var value = os.Getenv(key)
	if value == "" {
		return d
	}
	return value
}

func NewRepository() (*EventRepository, error) {
	connectionString := fmt.Sprintf(
		"user=MicroserviceApp dbname=BookingDB sslmode=disable password=MicroserviceApp host=%s port=%s",
		POSTGRES_HOST,
		POSTGRES_PORT,
	)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	db.MustExec(`
		DROP TABLE IF EXISTS booking_events;
		CREATE TABLE IF NOT EXISTS booking_events (
			sequence_nr SERIAL,
			entity_id uuid,
			type text,
			payload json,
			created_at timestamp default now()
	  );
	`)

	return &EventRepository{
		db: db,
	}, nil
}

func (r *EventRepository) Load(ctx context.Context, id string) (*ApartmentEntity, error) {

	records, err := r.db.Queryx(`SELECT * FROM booking_events WHERE
		 entity_id = $1 ORDER BY created_at`, id)

	if err != nil {
		return nil, err
	}

	var eventList []Event
	for records.Next() {
		var record = EventRecord{}
		err := records.StructScan(&record)
		if err != nil {
			return nil, err
		}
		fmt.Printf("READ RECORD: %+v\n", record)

		eventList = append(eventList, record.toEvent())
	}

	return NewFromEvents(eventList), nil
}

func (r *EventRepository) Save(ctx context.Context, b *ApartmentEntity) error {
	events := b.Events()
	for _, e := range events {
		fmt.Printf("EVENT: %+v\n", e)

		err := r.SaveEvent(ctx, b.ApartmentID, e)
		if err != nil {
			return err
		}
	}
	return nil
}

type EventRecord struct {
	SequenceNr int       `db:"sequence_nr"`
	EntityID   string    `db:"entity_id"`
	Type       string    `db:"type"`
	Payload    []byte    `db:"payload"`
	CreatedAt  time.Time `db:"created_at"`
}

func (r EventRecord) toEvent() Event {
	switch r.Type {
	case "ApartmentCreatedEvent":
		event := ApartmentCreatedEvent{}
		err := json.Unmarshal(r.Payload, &event)
		if err != nil {
			panic(err)
		}
		return event
	case "BookingCreatedEvent":
		event := BookingCreatedEvent{}
		err := json.Unmarshal(r.Payload, &event)
		if err != nil {
			panic(err)
		}
		return event
	case "BookingCancelledEvent":

		event := BookingCancelledEvent{}
		err := json.Unmarshal(r.Payload, &event)
		if err != nil {
			panic(err)
		}
		return event

	case "BookingUpdatedEvent":
		event := BookingUpdatedEvent{}
		err := json.Unmarshal(r.Payload, &event)
		if err != nil {
			panic(err)
		}
		return event
	}

	return nil
}

func (r *EventRepository) SaveEvent(ctx context.Context, entityID string, e Event) error {
	// Using reflection to get the concrete type
	value := reflect.ValueOf(e)

	payload, _ := json.Marshal(e)

	var record = EventRecord{
		EntityID:  entityID,
		Type:      value.Type().Name(),
		Payload:   payload,
		CreatedAt: time.Now(),
	}

	// fmt.Printf("WRITE RECORD: %+v\n", record)

	_, err := r.db.NamedExec(`INSERT INTO booking_events (
		entity_id,
		type,
		payload,
		created_at
	) 
		VALUES (
			:entity_id,
			:type,
			:payload,
			:created_at
		)`, &record)
	if err != nil {
		return err
	}
	return nil
}
