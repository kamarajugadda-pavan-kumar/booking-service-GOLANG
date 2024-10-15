package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/config"
	repository "github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/http/repositories"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/http/servicebase"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/types"
)
var cfg config.Config = config.MustGetConfig()
var servicebaseObj = servicebase.ServiceBase{BaseUrl: "http://localhost:3001"}

func BlockSeats(flightId string, numOfSeats int) error {
	type BlockFlightSeatsBody struct {
		NoOfSeats int    `json:"noOfSeats"`
		Action    string `json:"action"`
	}
	_, err := servicebaseObj.PUT_POST_PATCH(
		servicebase.MethodPatch,
		"/flight-search"+"/api/v1/flight/"+flightId,
		BlockFlightSeatsBody{NoOfSeats: numOfSeats, Action: "decrease"})
	if err != nil {
		fmt.Println("error in blocking seats", err)
		return err
	}
	return nil
}

func UnblockSeats(flightId string, numOfSeats int) error {
	type UnBlockFlightSeatsBody struct {
		NoOfSeats int    `json:"noOfSeats"`
		Action    string `json:"action"`
	}
	_, err := servicebaseObj.PUT_POST_PATCH(
		servicebase.MethodPatch,
		"/flight-search"+"/api/v1/flight/"+flightId,
		UnBlockFlightSeatsBody{NoOfSeats: numOfSeats, Action: "increase"})
	if err != nil {
		return err
	}
	return nil
}

func FetchFlightDetails(flightId string) (types.FlightData, error) {
	flightResponse, err := servicebaseObj.GET("/flight-search"+"/api/v1/flight/" + flightId)
	if err != nil {
		fmt.Printf("Failed to fetch flight details: %s\n", err)
	}
	var apiResponse types.ApiResponse
	json.Unmarshal(flightResponse, &apiResponse)
	return apiResponse.Data, err
}

func MakeBooking(flightId string, userId string, numOfSeats int) error {

	flightDetails, err := FetchFlightDetails(flightId)
	if err != nil {
		return err
	}

	blockingSeatsErr := BlockSeats(flightId, numOfSeats)
	if blockingSeatsErr != nil {
		return blockingSeatsErr
	}

	bookingData := types.Booking{
		FlightID:   flightId,
		UserID:     userId,
		Status:     types.Initiated,
		NumOfSeats: int64(numOfSeats),
		TotalCost:  float64(int64(numOfSeats) * flightDetails.Price),
	}
	_, bookingErr := repository.BookingRepository(bookingData)
	if bookingErr != nil {
		// If booking fails, unblock seats
		UnblockSeats(flightId, numOfSeats)
		return errors.New("booking failed, seats unblocked")
	}

	return nil
}

func MakePayment(bookingId string) error {
	// Payment gateway integration goes here
	// For demonstration purposes, we'll simulate payment
	booking, err := repository.FetchBooking(bookingId)
	if err != nil {
		return err
	}
	if booking.Status == types.Cancelled {
		return errors.New("booking is already cancelled")
	}
	if booking.Status == types.Booked {
		return errors.New("payment is already done")
	}
	// if booking.CreatedAt is past 10 minutes ago, cancel the booking
	createdAt, err := time.Parse(time.RFC3339, booking.CreatedAt)
	if err != nil {
		return errors.New("invalid booking creation time")
	}
	if createdAt.Add(10 * 60 * time.Second).Before(time.Now()) {
		return errors.New("booking expired due to inactivity")
	}
	_, paymentErr := repository.MakePayment(bookingId)
	if paymentErr != nil {
		return paymentErr
	}
	return nil
}
