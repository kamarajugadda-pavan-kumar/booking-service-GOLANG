package handlers_v1

import (
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/db"
	service "github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/http/services"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/types"
	utils "github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/utils/common"
)

func RegisterBookingRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/bookings/{flightID}", CreateBooking()).Methods("POST")
	router.HandleFunc("/api/v1/bookings/make-payment/{bookingID}", MakePayment()).Methods("PATCH")
	router.HandleFunc("/api/v1/bookings/cancel-booking/{bookingID}", CancelBooking()).Methods("PATCH")
	router.HandleFunc("/api/v1/bookings/{bookingID}", FetchBooking()).Methods("GET")
	router.HandleFunc("/api/v1/seats/{bookingId}", GetSeatsHandler).Methods("GET")
	router.HandleFunc("/api/v1/seat-booking", CreateSeatBooking).Methods("POST")
	router.HandleFunc("/api/v1/booking-history/{userId}", GetBookingHistory).Methods("GET")
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

		response := utils.Response{}

		numOfSeats, err := strconv.ParseInt(numOfSeatsStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid numOfSeats", http.StatusBadRequest)
			return
		}

		res, bookingError := service.MakeBooking(flightIDStr, userIdStr, int(numOfSeats))
		if bookingError != nil {
			http.Error(w, bookingError.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		successRes, err := response.SuccessResponse(
			res,
			"Booking was initialised, complete payment to confirm booking")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(successRes)

	}
}

func MakePayment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bookingIDStr := vars["bookingID"]

		paymentError := service.MakePayment(bookingIDStr)
		if paymentError != nil {
			http.Error(w, paymentError.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := utils.Response{}
		successRes, err := response.SuccessResponse(nil, "Payment successful, booking is confirmed")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(successRes)
	}
}

func CancelBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := utils.Response{}
		vars := mux.Vars(r)
		bookingIDStr := vars["bookingID"]

		bookingCancelResponse, cancelError := service.CancelBooking(bookingIDStr)
		if cancelError != nil {
			http.Error(w, cancelError.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		successRes, err := response.SuccessResponse(
			nil,
			bookingCancelResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(successRes)
	}
}

func FetchBooking() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := utils.Response{}

		vars := mux.Vars(r)
		bookingIDStr := vars["bookingID"]

		booking, err := service.FetchBooking(bookingIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		successRes, err := response.SuccessResponse(booking, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(successRes)

	}
}

func GetSeatsHandler(w http.ResponseWriter, r *http.Request) {
	database := db.GetDB()

	// Get bookingId from the URL path
	params := mux.Vars(r)
	bookingId := params["bookingId"]

	// Fetch flightId associated with the booking
	var flightId int
	err := database.QueryRow("SELECT flightId FROM Booking WHERE bookingId = ?", bookingId).Scan(&flightId)
	if err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	// Query for booked seats for this booking
	bookedSeats, err := service.GetBookedSeats(flightId, bookingId)
	if err != nil {
		http.Error(w, "Error fetching booked seats", http.StatusInternalServerError)
		return
	}

	// Query for available seats for this flight
	availableSeats, err := service.GetAvailableSeats(flightId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error fetching available seats", http.StatusInternalServerError)
		return
	}

	// Create response with both booked and available seats
	response := types.GetSeatsResponse{
		BookedSeats:    bookedSeats,
		AvailableSeats: availableSeats,
	}

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateSeatBooking(w http.ResponseWriter, r *http.Request) {
	var request types.SeatBookingRequest
	database := db.GetDB()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := database.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	stmt, err := tx.Prepare(`
        INSERT INTO SeatBookings (seatId, bookingId, flightId, name, age)
        VALUES (?, ?, ?, ?, ?)
    `)
	if err != nil {
		tx.Rollback()
		log.Printf("Error preparing statement: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, passenger := range request.Passengers {
		_, err := stmt.Exec(passenger.SeatID, request.BookingID, request.FlightID, passenger.Name, passenger.Age)
		if err != nil {
			tx.Rollback()
			log.Printf("Error inserting seat booking: %v", err)
			http.Error(w, fmt.Sprintf("Failed to book seat %d: %v", passenger.SeatID, err), http.StatusInternalServerError)
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Seat bookings created successfully"))
}

func GetBookingHistory(w http.ResponseWriter, r *http.Request) {

	userIdStr := mux.Vars(r)["userId"]

	bookings, err := service.FetchBookingHistory(userIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookings)
}
