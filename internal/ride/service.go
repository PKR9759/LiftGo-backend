// internal/ride/service.go
package ride

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, driverID string, req CreateRequest) (*Ride, error) {
	if req.OriginCity == "" || req.DestinationCity == "" {
		return nil, errors.New("origin and destination cities are required")
	}
	if req.OriginAddress == "" || req.DestinationAddress == "" {
		return nil, errors.New("origin and destination addresses are required")
	}
	if req.TotalSeats < 1 {
		return nil, errors.New("total seats must be at least 1")
	}
	if req.PricePerSeat < 0 {
		return nil, errors.New("price cannot be negative")
	}

	departure, err := time.Parse(time.RFC3339, req.DepartureAt)
	if err != nil {
		return nil, errors.New("invalid departure_at format, use ISO 8601")
	}
	if departure.Before(time.Now()) {
		return nil, errors.New("departure time must be in the future")
	}

	return s.repo.Create(ctx, driverID, req)
}

func (s *Service) GetByID(ctx context.Context, id string) (*Ride, error) {
	ride, err := s.repo.GetByID(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("ride not found")
	}
	return ride, err
}

func (s *Service) Search(ctx context.Context, p SearchParams) ([]*Ride, error) {
	if p.From == "" || p.To == "" {
		return nil, errors.New("from and to cities are required")
	}
	return s.repo.Search(ctx, p)
}

func (s *Service) GetMyRides(ctx context.Context, driverID string) ([]*Ride, error) {
	return s.repo.GetByDriver(ctx, driverID)
}

func (s *Service) Update(ctx context.Context, id, driverID string, req UpdateRequest) (*Ride, error) {
	if req.DepartureAt != "" {
		departure, err := time.Parse(time.RFC3339, req.DepartureAt)
		if err != nil {
			return nil, errors.New("invalid departure_at format, use ISO 8601")
		}
		if departure.Before(time.Now()) {
			return nil, errors.New("departure time must be in the future")
		}
	}
	ride, err := s.repo.Update(ctx, id, driverID, req)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.New("ride not found or you are not the driver")
	}
	return ride, err
}

func (s *Service) Cancel(ctx context.Context, id, driverID string) error {
	return s.repo.Cancel(ctx, id, driverID)
}