package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/db"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/types"
)

func BookingRepository(booking types.Booking) (types.BookingSucessData, error) {
	database := db.GetDB()
	successResponse := types.BookingSucessData{}

	tx, err := database.Begin()
	if err != nil {
		return successResponse, err
	}

	// Defer rollback in case anything fails
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `INSERT INTO Booking (userId, flightId, numOfSeats, totalCost, status) 
	          VALUES (?, ?, ?, ?, ?)`
	result, err := tx.Exec(query, booking.UserID, booking.FlightID, booking.NumOfSeats, booking.TotalCost, booking.Status)
	if err != nil {
		tx.Rollback()
		return successResponse, errors.New("booking failed: " + err.Error())
	}

	// Retrieve the last inserted ID
	insertedID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return successResponse, errors.New("failed to retrieve inserted ID: " + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return successResponse, errors.New("failed to commit transaction: " + err.Error())
	}
	successResponse.BookingID = int(insertedID)

	return successResponse, nil
}

func MakePayment(bookingId string) (string, error) {
	database := db.GetDB()

	tx, err := database.Begin()
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `UPDATE Booking SET (status) VALUES (?) WHERE id = ?`
	_, err = tx.Exec(query, types.Booked, bookingId)
	if err != nil {
		tx.Rollback() // Rollback the transaction if update fails
		return "", errors.New("payment failed: " + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return "", errors.New("failed to commit transaction: " + err.Error())
	}

	return fmt.Sprintf("Payment created successfully for booking ID: %s", bookingId), nil
}

func FetchBooking(bookingId string) (types.Booking, error) {
	database := db.GetDB()
	query := `SELECT * FROM Booking WHERE bookingId =?`
	var booking types.Booking
	row := database.QueryRow(query, bookingId)
	err := row.Scan(&booking.BookingID, &booking.FlightID, &booking.UserID, &booking.Status, &booking.NumOfSeats, &booking.TotalCost, &booking.CreatedAt, &booking.UpdatedAt)
	if err == sql.ErrNoRows {
		return booking, errors.New("booking not found")
	} else if err != nil {
		return booking, err
	}
	return booking, nil
}

func CancelBooking(bookingId string) (string, error) {
	database := db.GetDB()

	tx, err := database.Begin()
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `UPDATE Booking SET status = ? WHERE bookingId = ?`
	_, err = tx.Exec(query, types.Cancelled, bookingId)
	if err != nil {
		tx.Rollback() // Rollback the transaction if update fails
		return "", errors.New("cancellation failed: " + err.Error())
	}

	if err := tx.Commit(); err != nil {
		return "", errors.New("failed to commit transaction: " + err.Error())
	}

	return fmt.Sprintf("Booking cancelled successfully for booking ID: %s", bookingId), nil
}
