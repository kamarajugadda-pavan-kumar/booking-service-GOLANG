package types

type Status string

// Enum values for Status
const (
	Booked    Status = "booked"
	Cancelled Status = "cancelled"
	Initiated Status = "initiated"
	Pending   Status = "pending"
)

type Booking struct {
	FlightID   string `json:"flight_id"`
	UserID     string `json:"user_id"`
	Status     Status `json:"status"`
	NumOfSeats int64  `json:"num_of_seats"`
	TotalCost  int64  `json:"total_cost"`
}

type CreateBookingBody struct {
	UserID     int64 `json:"userId"`
	NumOfSeats int64 `json:"numOfSeats"`
}
