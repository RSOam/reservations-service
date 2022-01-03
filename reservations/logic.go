package reservations

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	consulapi "github.com/hashicorp/consul/api"
)

type service struct {
	db     ReservationDB
	logger log.Logger
	consul consulapi.Client
}

func NewService(db ReservationDB, logger log.Logger, consul consulapi.Client) ReservationsService {
	return &service{
		db:     db,
		logger: logger,
		consul: consul,
	}
}

func (s service) CreateReservation(ctx context.Context, from string, to string, userID string, chargerID string) (string, error) {
	logger := log.With(s.logger, "method: ", "CreateReservation")

	if err := s.db.CreateReservation(ctx, from, to, userID, chargerID); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	logger.Log("create Reservation", nil)
	return "Ok", nil
}
func (s service) GetReservation(ctx context.Context, id string) (Reservation, error) {
	logger := log.With(s.logger, "method", "GetReservation")
	reservation, err := s.db.GetReservation(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return reservation, err
	}
	logger.Log("Get Reservation", id)
	return reservation, nil
}
func (s service) GetReservations(ctx context.Context) ([]Reservation, error) {
	logger := log.With(s.logger, "method", "GetReservation")
	reservations, err := s.db.GetReservations(ctx)
	if err != nil {
		level.Error(logger).Log("err", err)
		return reservations, err
	}
	logger.Log("Get Reservations")
	return reservations, nil
}

func (s service) GetReservationsFilter(ctx context.Context, chargerID string, userID string) ([]Reservation, error) {
	logger := log.With(s.logger, "method", "GetReservationsFilter")
	reservations, err := s.db.GetReservationsFilter(ctx, chargerID, userID)
	if err != nil {
		level.Error(logger).Log("err", err)
		return reservations, err
	}
	logger.Log("Get ReservationsFilter")
	return reservations, nil
}
func (s service) DeleteReservation(ctx context.Context, id string) (string, error) {
	logger := log.With(s.logger, "method", "DeleteReservation")
	err := s.db.DeleteReservation(ctx, id)
	if err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	logger.Log("Delete Rating", id)
	return "Ok", nil
}
func (s service) UpdateReservation(ctx context.Context, id string, from string, to string) (string, error) {
	logger := log.With(s.logger, "method: ", "UpdateRating")

	if err := s.db.UpdateReservation(ctx, id, from, to); err != nil {
		level.Error(logger).Log("err", err)
		return "", err
	}
	logger.Log("update Rating", id)
	return "Ok", nil
}
