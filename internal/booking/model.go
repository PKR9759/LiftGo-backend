// internal/booking/model.go
package booking

import "time"


type Booking struct {
	ID              string    `json:"id"`
	RideID          string    `json:"ride_id"`
	RiderID         string    `json:"rider_id"`
	RiderName       string    `json:"rider_name"`
	DriverID        string    `json:"driver_id"`   
	DriverName      string    `json:"driver_name"`
	OriginCity      string    `json:"origin_city"`
	DestinationCity string    `json:"destination_city"`
	DepartureAt     time.Time `json:"departure_at"`
	Seats           int       `json:"seats"`
	Status          string    `json:"status"`
	TotalPrice      float64   `json:"total_price"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateRequest struct {
	RideID string `json:"ride_id"`
	Seats  int    `json:"seats"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}