package apartments

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Apartment struct {
	Id            string `db:"id" json:"id"`
	ApartmentName string `db:"apartment_name" json:"apartment_name"`
	Address       string `db:"address" json:"address"`
	NoiseLevel    string `db:"noise_level" json:"noise_level"`
	Floor         string `db:"floor" json:"floor"`
}

type ApartmentRepository struct {
	db *sqlx.DB
}

func NewApartmentRepository() (*ApartmentRepository, error) {
	var db, err = ConnectToDatabase()
	if err != nil {
		return nil, err
	}

	return &ApartmentRepository{
		db: db,
	}, nil
}

func ConnectToDatabase() (*sqlx.DB, error) {
	connectionString := fmt.Sprintf(
		"user=MicroserviceApp dbname=ApartmentDB sslmode=disable password=MicroserviceApp host=%s port=%s",
		POSTGRES_HOST,
		POSTGRES_PORT,
	)

	db, err := sqlx.Connect("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	var schema = `
	DROP TABLE IF EXISTS apartments;

	CREATE TABLE IF NOT EXISTS apartments (
		id uuid primary key DEFAULT gen_random_uuid(),
		apartment_name text,
		address text,
		noise_level text,
		floor text
	);`
	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO apartments (apartment_name, address, noise_level, floor) VALUES ($1, $2, $3, $4)", "Always Green", "Bolzano", "2", "3")
	tx.MustExec("INSERT INTO apartments (apartment_name, address, noise_level, floor) VALUES ($1, $2, $3, $4)", "Rarely Yellow", "Bolzano", "4", "3")
	tx.NamedExec("INSERT INTO apartments (apartment_name, address, noise_level, floor) VALUES (:apartment_name, :address, :noise_level, :floor)", &Apartment{"0", "Sometimes Pink", "Merano", "1", "5"})
	tx.Commit()

	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Successfully setup database connection")
	return db, nil
}

func (a *ApartmentRepository) SaveApartment(apartment Apartment) Apartment {
	apartment.Id = uuid.NewString()
	_, err := a.db.NamedExec("INSERT INTO apartments (id, apartment_name, address, noise_level, floor) VALUES (:id, :apartment_name, :address, :noise_level, :floor)", &apartment)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Apartment added: %v\n", apartment)
	return apartment
}

func (a *ApartmentRepository) DeleteApartment(apartmentId string) {
	_, err := a.db.Exec("DELETE FROM apartments WHERE id = $1", apartmentId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Deleted apartment with id: %v\n", apartmentId)
}

func (a *ApartmentRepository) ListAllApartments() []Apartment {
	apartment := Apartment{}
	var apartmentList []Apartment

	rows, _ := a.db.Queryx("SELECT * FROM apartments")

	for rows.Next() {
		err := rows.StructScan(&apartment)
		if err != nil {
			log.Fatalln(err)
		}
		apartmentList = append(apartmentList, apartment)
	}
	return apartmentList
}
