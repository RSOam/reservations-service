package reservations

import (
	"context"
)

type ReservationsService interface {
	CreateReservation(ctx context.Context, from string, to string, userToken string, chargerID string) (string, error)
	GetReservation(ctx context.Context, id string) (Reservation, error)
	GetReservations(ctx context.Context) ([]Reservation, error)
	GetReservationsFilter(ctx context.Context, chargerID string, userID string) ([]Reservation, error)
	UpdateReservation(ctx context.Context, id string, from string, to string) (string, error)
	DeleteReservation(ctx context.Context, id string) (string, error)
	ReservationClosest(ctx context.Context, userToken string, from string, to string, location Location) (Reservation, string, error)
}
