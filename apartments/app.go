package apartments

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func apartmentsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		fmt.Printf("got /api/apartments GET request\n")
		w.Header().Set("Content-Type", "application/json")
		allApartments := ListAllApartments()
		json.NewEncoder(w).Encode(&allApartments)
	case http.MethodPost:
		var apartment Apartment
		err := json.NewDecoder(r.Body).Decode(&apartment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Printf("got /api/apartments POST request\n")
		SaveApartment(apartment)
		json.NewEncoder(w).Encode(&apartment)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func StartApp() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/api/apartments", apartmentsHandler)

	err := http.ListenAndServe(":3000", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
