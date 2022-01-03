package reservations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Charger struct {
	ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name          string             `json:"name"`
	Location      Location           `json:"location"`
	AverageRating float64            `json:"averageRating"`
	Ratings       []Rating           `json:"ratings"`
	Comments      []Comment          `json:"comments"`
	Reservations  []Reservation      `json:"reservations"`
	Created       string             `json:"created"`
	Modified      string             `json:"modified"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Comment struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChargerID primitive.ObjectID `json:"chargerID"`
	UserID    primitive.ObjectID `json:"userID"`
	Text      string             `json:"text"`
	Created   string             `json:"created"`
	Modified  string             `json:"modified"`
}
type Rating struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChargerID primitive.ObjectID `json:"chargerID"`
	UserID    primitive.ObjectID `json:"userID"`
	Rating    int                `json:"rating"`
	Created   string             `json:"created"`
	Modified  string             `json:"modified"`
}

type Reservation struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChargerID primitive.ObjectID `json:"chargerID"`
	UserID    primitive.ObjectID `json:"userID"`
	From      string             `json:"from"`
	To        string             `json:"to"`
	Created   string             `json:"created"`
	Modified  string             `json:"modified"`
}

type ReservationDB interface {
	CreateReservation(ctx context.Context, from string, to string, userID string, chargerID string) error
	GetReservation(ctx context.Context, id string) (Reservation, error)
	GetReservations(ctx context.Context) ([]Reservation, error)
	GetReservationsFilter(ctx context.Context, chargerID string, userID string) ([]Reservation, error)
	UpdateReservation(ctx context.Context, id string, from string, to string) error
	DeleteReservation(ctx context.Context, id string) error
	ReservationClosest(ctx context.Context, userID, from string, to string, location Location) (Reservation, error)
}
