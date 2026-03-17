// internal/ride/repository.go
package ride

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, driverID string, req CreateRequest) (*Ride, error) {
	departure, err := time.Parse(time.RFC3339, req.DepartureAt)
	if err != nil {
		return nil, fmt.Errorf("invalid departure_at format, use ISO 8601")
	}

	ride := &Ride{}
	err = r.db.QueryRow(ctx,
		`INSERT INTO rides
			(driver_id, origin_city, destination_city, origin_address,
			 destination_address, departure_at, total_seats, available_seats,
			 price_per_seat, notes)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$7,$8,$9)
		 RETURNING id, driver_id, origin_city, destination_city,
		           origin_address, destination_address, departure_at,
		           total_seats, available_seats, price_per_seat,
		           notes, status, created_at`,
		driverID,
		req.OriginCity, req.DestinationCity,
		req.OriginAddress, req.DestinationAddress,
		departure, req.TotalSeats,
		req.PricePerSeat, nullableString(req.Notes),
	).Scan(
		&ride.ID, &ride.DriverID,
		&ride.OriginCity, &ride.DestinationCity,
		&ride.OriginAddress, &ride.DestinationAddress,
		&ride.DepartureAt, &ride.TotalSeats, &ride.AvailableSeats,
		&ride.PricePerSeat, &ride.Notes, &ride.Status, &ride.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ride, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Ride, error) {
	ride := &Ride{}
	err := r.db.QueryRow(ctx,
		`SELECT r.id, r.driver_id, u.name, u.avg_rating, u.total_reviews,
		        r.origin_city, r.destination_city, r.origin_address,
		        r.destination_address, r.departure_at, r.total_seats,
		        r.available_seats, r.price_per_seat, r.notes,
		        r.status, r.created_at
		 FROM rides r
		 JOIN users u ON u.id = r.driver_id
		 WHERE r.id = $1`, id,
	).Scan(
		&ride.ID, &ride.DriverID, &ride.DriverName,
		&ride.DriverAvgRating, &ride.DriverTotalReviews,
		&ride.OriginCity, &ride.DestinationCity,
		&ride.OriginAddress, &ride.DestinationAddress,
		&ride.DepartureAt, &ride.TotalSeats, &ride.AvailableSeats,
		&ride.PricePerSeat, &ride.Notes, &ride.Status, &ride.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ride, nil
}

func (r *Repository) Search(ctx context.Context, p SearchParams) ([]*Ride, error) {
	seats := p.Seats
	if seats < 1 {
		seats = 1
	}

	// build date range: full day from 00:00 to 23:59
	var dayStart, dayEnd time.Time
	if p.Date != "" {
		var err error
		dayStart, err = time.Parse("2006-01-02", p.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD")
		}
		dayEnd = dayStart.Add(24*time.Hour - time.Second)
	} else {
		dayStart = time.Now()
		dayEnd = dayStart.Add(365 * 24 * time.Hour)
	}

	rows, err := r.db.Query(ctx,
		`SELECT r.id, r.driver_id, u.name, u.avg_rating, u.total_reviews,
		        r.origin_city, r.destination_city, r.origin_address,
		        r.destination_address, r.departure_at, r.total_seats,
		        r.available_seats, r.price_per_seat, r.notes,
		        r.status, r.created_at
		 FROM rides r
		 JOIN users u ON u.id = r.driver_id
		 WHERE r.status = 'active'
		   AND r.available_seats >= $1
		   AND LOWER(r.origin_city) LIKE LOWER($2)
		   AND LOWER(r.destination_city) LIKE LOWER($3)
		   AND r.departure_at BETWEEN $4 AND $5
		 ORDER BY r.departure_at ASC`,
		seats,
		"%"+p.From+"%",
		"%"+p.To+"%",
		dayStart, dayEnd,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []*Ride
	for rows.Next() {
		ride := &Ride{}
		err := rows.Scan(
			&ride.ID, &ride.DriverID, &ride.DriverName,
			&ride.DriverAvgRating, &ride.DriverTotalReviews,
			&ride.OriginCity, &ride.DestinationCity,
			&ride.OriginAddress, &ride.DestinationAddress,
			&ride.DepartureAt, &ride.TotalSeats, &ride.AvailableSeats,
			&ride.PricePerSeat, &ride.Notes, &ride.Status, &ride.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}
	return rides, nil
}

func (r *Repository) GetByDriver(ctx context.Context, driverID string) ([]*Ride, error) {
	rows, err := r.db.Query(ctx,
		`SELECT r.id, r.driver_id, u.name, u.avg_rating, u.total_reviews,
		        r.origin_city, r.destination_city, r.origin_address,
		        r.destination_address, r.departure_at, r.total_seats,
		        r.available_seats, r.price_per_seat, r.notes,
		        r.status, r.created_at
		 FROM rides r
		 JOIN users u ON u.id = r.driver_id
		 WHERE r.driver_id = $1
		 ORDER BY r.departure_at DESC`, driverID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rides []*Ride
	for rows.Next() {
		ride := &Ride{}
		err := rows.Scan(
			&ride.ID, &ride.DriverID, &ride.DriverName,
			&ride.DriverAvgRating, &ride.DriverTotalReviews,
			&ride.OriginCity, &ride.DestinationCity,
			&ride.OriginAddress, &ride.DestinationAddress,
			&ride.DepartureAt, &ride.TotalSeats, &ride.AvailableSeats,
			&ride.PricePerSeat, &ride.Notes, &ride.Status, &ride.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		rides = append(rides, ride)
	}
	return rides, nil
}

func (r *Repository) Update(ctx context.Context, id, driverID string, req UpdateRequest) (*Ride, error) {
	departure, err := time.Parse(time.RFC3339, req.DepartureAt)
	if err != nil {
		return nil, fmt.Errorf("invalid departure_at format, use ISO 8601")
	}

	ride := &Ride{}
	err = r.db.QueryRow(ctx,
		`UPDATE rides
		 SET origin_address      = COALESCE(NULLIF($1,''), origin_address),
		     destination_address = COALESCE(NULLIF($2,''), destination_address),
		     departure_at        = $3,
		     total_seats         = CASE WHEN $4 > 0 THEN $4 ELSE total_seats END,
		     price_per_seat      = CASE WHEN $5 > 0 THEN $5 ELSE price_per_seat END,
		     notes               = COALESCE(NULLIF($6,''), notes),
		     updated_at          = now()
		 WHERE id = $7 AND driver_id = $8
		 RETURNING id, driver_id, origin_city, destination_city,
		           origin_address, destination_address, departure_at,
		           total_seats, available_seats, price_per_seat,
		           notes, status, created_at`,
		req.OriginAddress, req.DestinationAddress,
		departure, req.TotalSeats, req.PricePerSeat,
		nullableString(req.Notes), id, driverID,
	).Scan(
		&ride.ID, &ride.DriverID,
		&ride.OriginCity, &ride.DestinationCity,
		&ride.OriginAddress, &ride.DestinationAddress,
		&ride.DepartureAt, &ride.TotalSeats, &ride.AvailableSeats,
		&ride.PricePerSeat, &ride.Notes, &ride.Status, &ride.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return ride, nil
}

func (r *Repository) Cancel(ctx context.Context, id, driverID string) error {
	result, err := r.db.Exec(ctx,
		`UPDATE rides SET status = 'cancelled', updated_at = now()
		 WHERE id = $1 AND driver_id = $2 AND status = 'active'`,
		id, driverID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("ride not found or already cancelled")
	}
	return nil
}

// nullableString converts empty string to nil so DB stores NULL
func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}