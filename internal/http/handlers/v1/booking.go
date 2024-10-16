package handlers_v1

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	service "github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/http/services"
	utils "github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/utils/common"
)

func RegisterBookingRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/bookings/{flightID}", CreateBooking()).Methods("POST")
	router.HandleFunc("/api/v1/bookings/make-payment/{bookingID}", MakePayment()).Methods("PATCH")
	router.HandleFunc("/api/v1/bookings/cancel-booking/{bookingID}", CancelBooking()).Methods("PATCH")
	router.HandleFunc("/api/v1/bookings/{bookingID}", FetchBooking()).Methods("GET")
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
