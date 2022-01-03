package reservations

import (
	"context"
	"net/http"

	ht "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewHttpServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/reservations").Handler(ht.NewServer(
		endpoints.CreateReservation,
		decodeCreateReservationRequest,
		encodeResponse,
	))
	r.Methods("PUT").Path("/reservations/{id}").Handler(ht.NewServer(
		endpoints.UpdateReservation,
		decodeUpdateReservationRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/reservations/{id}").Handler(ht.NewServer(
		endpoints.GetReservation,
		decodeGetReservationRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/reservations").Handler(ht.NewServer(
		endpoints.GetReservations,
		decodeGetReservationsRequest,
		encodeResponse,
	))
	r.Methods("GET").Path("/reservations/").Handler(ht.NewServer(
		endpoints.GetReservationsFilter,
		decodeGetReservationsFilterRequest,
		encodeResponse,
	))
	r.Methods("DELETE").Path("/reservations/{id}").Handler(ht.NewServer(
		endpoints.DeleteReservation,
		decodeDeleteReservationRequest,
		encodeResponse,
	))
	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
