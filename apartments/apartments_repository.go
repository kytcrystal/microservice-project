package apartments

import "fmt"

var ApartmentList []Apartment

func SaveApartment(apartment Apartment) {
	ApartmentList = append(ApartmentList, apartment)
	fmt.Printf("Apartments: %v\n", ApartmentList)
}

func ListAllApartments() []Apartment {
	return ApartmentList
}
