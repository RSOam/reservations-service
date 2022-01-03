package reservations

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateReservation     endpoint.Endpoint
	GetReservation        endpoint.Endpoint
	GetReservations       endpoint.Endpoint
	GetReservationsFilter endpoint.Endpoint
	UpdateReservation     endpoint.Endpoint
	DeleteReservation     endpoint.Endpoint
	ReservationClosest    endpoint.Endpoint
}

func MakeEndpoints(s ReservationsService) Endpoints {
	return Endpoints{
		CreateReservation:     makeCreateReservationEndpoint(s),
		GetReservation:        makeGetReservationEndpoint(s),
		GetReservations:       makeGetReservationsEndpoint(s),
		GetReservationsFilter: makeGetReservationsFilterEndpoint(s),
		UpdateReservation:     makeUpdateReservationEndpoint(s),
		DeleteReservation:     makeDeleteReservationEndpoint(s),
		ReservationClosest:    makeReservationClosestEndpoint(s),
	}
}

func makeCreateReservationEndpoint(s ReservationsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateReservationRequest)
		status, err := s.CreateReservation(ctx, req.From, req.To, req.UserToken, req.ChargerID)
		return CreateReservationResponse{Status: status}, err
	}
}
func makeGetReservationEndpoint(s ReservationsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReservationRequest)
		reservation, err := s.GetReservation(ctx, req.Id)
		return GetReservationResponse{
			ChargerID: reservation.ChargerID.Hex(),
			UserID:    reservation.UserID.Hex(),
			From:      reservation.From,
			To:        reservation.To,
			Created:   reservation.Created,
			Modified:  reservation.Modified,
		}, err
	}
}
func makeGetReservationsEndpoint(s ReservationsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reservations, err := s.GetReservations(ctx)
		return GetReservationsResponse{
			Reservations: reservations,
		}, err
	}
}
func makeDeleteReservationEndpoint(s ReservationsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteReservationRequest)
		status, err := s.DeleteReservation(ctx, req.Id)
		return DeleteReservationResponse{
			Status: status,
		}, err
	}
}

func makeGetRatingsFilterEndpoint(s ReservationsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReservationsFilterRequest)
		reservations, err := s.GetReservationsFilter(ctx, req.ChargerID, req.UserID)
		return GetReservationsFilterResponse{
			Reservations: reservations,
		}, err
	}
}
func makeUpdateReservationEndpoint(s ReservationsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReservationRequest)
		status, err := s.UpdateReservation(ctx, req.Id, req.From, req.To)
		return CreateReservationResponse{Status: status}, err
	}
}
func makeGetReservationsFilterEndpoint(s ReservationsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReservationsFilterRequest)
		reservations, err := s.GetReservationsFilter(ctx, req.ChargerID, req.UserID)
		return GetReservationsFilterResponse{
			Reservations: reservations,
		}, err
	}
}
func makeReservationClosestEndpoint(s ReservationsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ReservationClosestRequest)
		reservation, status, err := s.ReservationClosest(ctx, req.UserToken, req.From, req.To, req.Location)
		if err != nil {
			return ReservationClosestResponse{
				ChargerID: reservation.ChargerID.Hex(),
				UserID:    reservation.UserID.Hex(),
				From:      reservation.From,
				To:        reservation.To,
				Created:   reservation.Created,
				Modified:  reservation.Modified,
			}, err
		}
		return CreateReservationResponse{Status: status}, err
	}
}
