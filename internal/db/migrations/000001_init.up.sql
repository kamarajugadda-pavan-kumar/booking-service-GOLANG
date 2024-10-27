CREATE TABLE
    IF NOT EXISTS Booking (
        bookingId INT AUTO_INCREMENT PRIMARY KEY,
        flightId INT NOT NULL,
        userId INT NOT NULL,
        status ENUM ('booked', 'cancelled', 'initiated', 'pending') NOT NULL,
        numOfSeats INT DEFAULT 1,
        totalCost DECIMAL(10, 2) NOT NULL,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (flightId) REFERENCES Flights (id) ON DELETE CASCADE
    );