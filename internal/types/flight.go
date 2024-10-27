package types

type FlightData struct {
	ID               int64  `json:"id"`
	FlightNumber     string `json:"flightNumber"`
	AirplaneID       int64  `json:"airplaneId"`
	DepartureAirport string `json:"departureAirportId"`
	ArrivalAirport   string `json:"arrivalAirportId"`
	ArrivalTime      string `json:"arrivalTime"`
	DepartureTime    string `json:"departureTime"`
	Price            int64  `json:"price"`
	BoardingGate     string `json:"boardingGate"`
	TotalSeats       int64  `json:"totalSeats"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type ApiResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    FlightData `json:"data"`
	Error   struct{}   `json:"error"`
}
