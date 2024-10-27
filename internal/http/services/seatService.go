package service

import (
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/db"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/types"
)

// getBookedSeats fetches the booked seats based on bookingId and flightId
func GetBookedSeats(flightId int, bookingId string) ([]types.Seat, error) {
	database := db.GetDB()
	query := `
		SELECT s.id, s.row, s.col, s.type 
		FROM Seats s 
		JOIN SeatBookings sb ON s.id = sb.seatId 
		WHERE sb.flightId = ?`
	rows, err := database.Query(query, flightId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []types.Seat
	for rows.Next() {
		var seat types.Seat
		if err := rows.Scan(&seat.SeatID, &seat.Row, &seat.Col, &seat.SeatType); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	return seats, nil
}

// getAvailableSeats fetches the available seats for a flightId
func GetAvailableSeats(flightId int) ([]types.Seat, error) {
	database := db.GetDB()
	query := `
		SELECT s.id, s.row, s.col, s.type 
		FROM Seats s
		LEFT JOIN SeatBookings sb ON s.id = sb.seatId
		LEFT JOIN Flights f ON f.airplaneId = s.airplaneId
		WHERE sb.seatId IS NULL
		AND f.id = ?`
	rows, err := database.Query(query, flightId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []types.Seat
	for rows.Next() {
		var seat types.Seat
		if err := rows.Scan(&seat.SeatID, &seat.Row, &seat.Col, &seat.SeatType); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	return seats, nil
}
