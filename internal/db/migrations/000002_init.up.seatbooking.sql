CREATE TABLE
    SeatBookings (
        id INT AUTO_INCREMENT PRIMARY KEY,
        seatId INT NOT NULL,
        bookingId INT NOT NULL,
        flightId INT NOT NULL,
        name VARCHAR(255) NOT NULL,
        age INT NOT NULL,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (seatId) REFERENCES Seats (id),
        FOREIGN KEY (bookingId) REFERENCES Booking (bookingId) ON DELETE CASCADE,
        FOREIGN KEY (flightId) REFERENCES Flights (id),
        CONSTRAINT unique_flight_booking_seat UNIQUE (flightId, bookingId, seatId)
    );