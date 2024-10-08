package types

type Flight struct {
	id                 int64  `json:"id"`
	flightNumber       string `json:"flightNumber"`
	airplaneId         int64  `json:"airplaneId"`
	departureAirportId string `json:"departureAirportId"`
	arrivalAirportId   string `json:"arrivalAirportId"`
	arrivalTime        string `json:"arrivalTime"`
	departureTime      string `json:"departureTime"`
	price              int64  `json:"price"`
	boardingGate       string `json:"boardingGate"`
	totalSeats         int64  `json:"totalSeats"`
	createdAt          string `json:"createdAt"`
	updatedAt          string `json:"updatedAt"`
}
