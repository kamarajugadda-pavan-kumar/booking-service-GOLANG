package handlers_v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterBookingRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/bookings/{flightID}", CreateBooking()).Methods("POST")
}

type CreateBookingBody struct {
	UserID     int64 `json:"userId"`
	NumOfSeats int64 `json:"numOfSeats"`
}

func CreateBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		flightIDStr := vars["flightID"]

		var body CreateBookingBody
		json.NewDecoder(r.Body).Decode(&body)

		// send response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Creating booking for flight ID: " + flightIDStr))

		// Additional business logic can go here

	}
}
