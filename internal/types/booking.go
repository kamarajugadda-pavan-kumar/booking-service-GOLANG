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
	BookingID  string  `json:"bookingId"`
	FlightID   string  `json:"flightId"`
	UserID     string  `json:"userId"`
	Status     Status  `json:"status"`
	NumOfSeats int64   `json:"numOfSeats"`
	TotalCost  float64 `json:"totalCost"`
	CreatedAt  string  `json:"createdAt"`
	UpdatedAt  string  `json:"updatedAt"`
}

type CreateBookingBody struct {
	UserID     int64 `json:"userId"`
	NumOfSeats int64 `json:"numOfSeats"`
}

type BookingSucessData struct {
	BookingID int `json:"bookingId"`
}
