package reservations

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type (
	CreateReservationRequest struct {
		ChargerID string `json:"chargerID"`
		UserID    string `json:"userID"`
		From      string `json:"from"`
		To        string `json:"to"`
	}
	CreateReservationResponse struct {
		Status string `json:"status"`
	}
	GetReservationRequest struct {
		Id string `json:"id"`
	}
	GetReservationResponse struct {
		ChargerID string `json:"chargerID"`
		UserID    string `json:"userID"`
		From      string `json:"from"`
		To        string `json:"to"`
		Created   string `json:"created"`
		Modified  string `json:"modified"`
	}
	GetReservationsRequest struct {
	}
	GetReservationsResponse struct {
		Reservations []Reservation `json:"reservations"`
	}
	UpdateReservationRequest struct {
		Id   string `json:"id"`
		From string `json:"from"`
		To   string `json:"to"`
	}
	UpdateReservationResponse struct {
		Status string `json:"status"`
	}
	DeleteReservationRequest struct {
		Id string `json:"id"`
	}
	DeleteReservationResponse struct {
		Status string `json:"status"`
	}
	GetReservationsFilterRequest struct {
		ChargerID string `json:"chargerID"`
		UserID    string `json:"userID"`
	}
	GetReservationsFilterResponse struct {
		Reservations []Reservation `json:"reservations"`
	}
	//OTHER
	GetChargerRatingsRequest struct {
	}
	GetChargerRatingsResponse struct {
		Ratings []Rating `json:"ratings"`
	}
	GetChargerCommentsRequest struct {
	}
	GetChargerCommentsResponse struct {
		Comments []Comment `json:"comments"`
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeCreateReservationRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := CreateReservationRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
func decodeUpdateReservationRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := UpdateReservationRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	vals := mux.Vars(r)
	req.Id = vals["id"]
	return req, nil
}
func decodeGetReservationRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := GetReservationRequest{}
	vals := mux.Vars(r)
	req.Id = vals["id"]
	return req, nil
}
func decodeGetReservationsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := GetReservationRequest{}
	return req, nil
}
func decodeDeleteReservationRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := DeleteReservationRequest{}
	vals := mux.Vars(r)
	req.Id = vals["id"]
	return req, nil
}
func decodeGetReservationsFilterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	req := GetReservationsFilterRequest{}
	req.ChargerID = r.URL.Query().Get("charger")
	req.UserID = r.URL.Query().Get("user")
	return req, nil
}
