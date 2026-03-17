// internal/ride/model.go
package ride

import "time"

type Ride struct {
	ID                 string    `json:"id"`
	DriverID           string    `json:"driver_id"`
	DriverName         string    `json:"driver_name"`
	DriverAvgRating    float64   `json:"driver_avg_rating"`
	DriverTotalReviews int       `json:"driver_total_reviews"`
	OriginCity         string    `json:"origin_city"`
	DestinationCity    string    `json:"destination_city"`
	OriginAddress      string    `json:"origin_address"`
	DestinationAddress string    `json:"destination_address"`
	DepartureAt        time.Time `json:"departure_at"`
	TotalSeats         int       `json:"total_seats"`
	AvailableSeats     int       `json:"available_seats"`
	PricePerSeat       float64   `json:"price_per_seat"`
	Notes              *string   `json:"notes,omitempty"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
}

type CreateRequest struct {
	OriginCity         string  `json:"origin_city"`
	DestinationCity    string  `json:"destination_city"`
	OriginAddress      string  `json:"origin_address"`
	DestinationAddress string  `json:"destination_address"`
	DepartureAt        string  `json:"departure_at"` // ISO 8601 string from frontend
	TotalSeats         int     `json:"total_seats"`
	PricePerSeat       float64 `json:"price_per_seat"`
	Notes              string  `json:"notes"`
}

type UpdateRequest struct {
	OriginAddress      string  `json:"origin_address"`
	DestinationAddress string  `json:"destination_address"`
	DepartureAt        string  `json:"departure_at"`
	TotalSeats         int     `json:"total_seats"`
	PricePerSeat       float64 `json:"price_per_seat"`
	Notes              string  `json:"notes"`
}

type SearchParams struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Date  string `json:"date"`  // YYYY-MM-DD
	Seats int    `json:"seats"`
}