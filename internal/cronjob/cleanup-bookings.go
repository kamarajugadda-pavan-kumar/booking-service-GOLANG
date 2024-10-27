package cronjob

import (
	"fmt"
	"time"

	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/db"
	service "github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/http/services"
	"github.com/kamarajugadda-pavan-kumar/booking-service-GOLANG/internal/types"
)

func CleanExpiredBookings() {
	ticker := time.NewTicker(10 * time.Minute) // Run every 10 minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			removeExpiredBookings()
		}
	}
}

func removeExpiredBookings() {
	fmt.Println("Cleaning expired bookings...")
	database := db.GetDB()

	results, err := database.Query(`
					SELECT bookingId, flightId, numOfSeats 
					FROM Booking 
					WHERE createdAt <= NOW() - INTERVAL 10 MINUTE 
					AND status != 'booked';`)
	if err != nil {
		fmt.Println("Error fetching expired bookings:", err)
		return
	}
	defer results.Close()

	for results.Next() {
		var booking types.Booking
		err := results.Scan(
			&booking.BookingID,
			&booking.FlightID,
			&booking.NumOfSeats,
		)
		if err != nil {
			fmt.Println("Error scanning expired booking:", err)
			continue
		}

		// Unblock the seats for the expired booking
		err = service.UnblockSeats(booking.FlightID, int(booking.NumOfSeats))
		if err != nil {
			fmt.Println("Error unblocking seats for expired booking:", err)
			continue
		}

		fmt.Printf("Seats unblocked for expired booking with ID %s\n", booking.BookingID)

		_, err = database.Exec("DELETE FROM Booking WHERE bookingId =?", booking.BookingID)
		if err != nil {
			fmt.Println("Error deleting expired booking:", err)
			continue
		}

		fmt.Printf("Expired booking with ID %s deleted\n", booking.BookingID)
	}
}
