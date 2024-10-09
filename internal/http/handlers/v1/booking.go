package handlers_v1

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	service "github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/http/services"
	types "github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/types"
)

func RegisterBookingRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/bookings/{flightID}", CreateBooking()).Methods("POST")
}

func CreateBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		flightIDStr := vars["flightID"]

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Extract form values
		userIdStr := r.FormValue("userId")
		numOfSeatsStr := r.FormValue("numOfSeats")

		// Convert form values to the appropriate types
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid userId", http.StatusBadRequest)
			return
		}
		numOfSeats, err := strconv.ParseInt(numOfSeatsStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid numOfSeats", http.StatusBadRequest)
			return
		}

		body := types.CreateBookingBody{UserID: userId, NumOfSeats: numOfSeats}

		bookingRes, err := service.CreateBookingService(flightIDStr, body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(bookingRes))

		// Additional business logic can go here

	}
}
