package repository

import (
	"errors"
	"fmt"

	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/db"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/types"
)

func BookingRepository(booking types.Booking) (string, error) {
	database := db.GetDB()
	// Start a transaction
	tx, err := database.Begin()
	if err != nil {
		return "", err
	}

	// Defer rollback in case anything fails
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Create booking query
	query := `INSERT INTO Booking (userId, flightId, numOfSeats, totalCost, status) 
	          VALUES (?, ?, ?, ?, ?)`
	_, err = tx.Exec(query, booking.UserID, booking.FlightID, booking.NumOfSeats, booking.TotalCost, booking.Status)
	if err != nil {
		tx.Rollback() // Rollback the transaction if insert fails
		return "", errors.New("booking failed: " + err.Error())
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return "", errors.New("failed to commit transaction: " + err.Error())
	}

	return fmt.Sprintf("Booking created successfully for user ID: %d", booking.UserID), nil

}
