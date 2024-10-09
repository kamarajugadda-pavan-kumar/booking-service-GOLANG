package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	repository "github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/http/repositories"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/types"
)

func CreateBookingService(flightID string, body types.CreateBookingBody) (string, error) {

	json.Marshal(body)

	response, err := http.NewRequest(http.MethodPatch,
		"http://localhost:3001/api/v1/flight/"+flightID,
		bytes.NewReader())

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	jsonResponse, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err)
		return "", errors.New("failed to read response body")
	}

	var apiResponse types.ApiResponse
	err1 := json.Unmarshal([]byte(jsonResponse), &apiResponse)
	if err1 != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return "", errors.New("failed to unmarshal response ")
	}

	fmt.Println(apiResponse.Data)

	var flightData *types.FlightData = &apiResponse.Data
	if flightData.TotalSeats < body.NumOfSeats {
		return "", errors.New("not enough seats available")
	}

	var totalCost int64 = flightData.Price * body.NumOfSeats

	// Create the booking in the database
	booking := types.Booking{
		FlightID:   2,
		UserID:     body.UserID,
		Status:     types.Initiated,
		NumOfSeats: body.NumOfSeats,
		TotalCost:  totalCost,
	}

	res, err := repository.BookingRepository(booking)
	if err != nil {
		return "", err
	}

	return res, nil

}
