package service

import (
	"encoding/json"
	"errors"
	"fmt"

	repository "github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/http/repositories"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/http/servicebase"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/types"
)

var servicebaseObj = servicebase.ServiceBase{BaseUrl: "http://localhost:3001"}

func BlockSeats(flightId string, numOfSeats int) error {
	type BlockFlightSeatsBody struct {
		NoOfSeats int    `json:"noOfSeats"`
		Action    string `json:"action"`
	}
	_, err := servicebaseObj.PUT_POST_PATCH(
		servicebase.MethodPatch,
		"/api/v1/flight/"+flightId,
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
		"/api/v1/flight/"+flightId,
		UnBlockFlightSeatsBody{NoOfSeats: numOfSeats, Action: "increase"})
	if err != nil {
		return err
	}
	return nil
}

func FetchFlightDetails(flightId string) (types.FlightData, error) {
	flightResponse, err := servicebaseObj.GET("/api/v1/flight/" + flightId)
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
