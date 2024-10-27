package types

type Seat struct {
	SeatID   int    `json:"seatId"`
	Row      int    `json:"row"`
	Col      string `json:"col"`
	SeatType string `json:"seatType"`
}

type GetSeatsResponse struct {
	BookedSeats    []Seat `json:"bookedSeats"`
	AvailableSeats []Seat `json:"availableSeats"`
}

type SeatBookingRequest struct {
	FlightID   int         `json:"flightId"`
	BookingID  string      `json:"bookingId"`
	Passengers []Passenger `json:"passengers"`
}

type Passenger struct {
	SeatID int    `json:"seatId"`
	Name   string `json:"name"`
	Age    string `json:"age"`
}
