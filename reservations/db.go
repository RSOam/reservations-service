package reservations

import (
	"context"
	"time"

	"github.com/go-kit/log"
	consulapi "github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type database struct {
	db     *mongo.Database
	logger log.Logger
	consul consulapi.Client
}

func NewDatabase(db *mongo.Database, logger log.Logger, consul consulapi.Client) ReservationDB {
	return &database{
		db:     db,
		logger: log.With(logger, "database", "mongoDB"),
		consul: consul,
	}
}

func (dat *database) CreateReservation(ctx context.Context, from string, to string, userID string, chargerID string) error {
	userIDmongo, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		dat.logger.Log("Error creating reservation: ", err.Error())
		return err
	}
	chargerIDmongo, err := primitive.ObjectIDFromHex(chargerID)
	if err != nil {
		dat.logger.Log("Error creating rating: ", err.Error())
		return err
	}
	reservationObj := Reservation{
		ChargerID: chargerIDmongo,
		UserID:    userIDmongo,
		From:      from,
		To:        to,
		Created:   time.Now().Format(time.RFC3339),
		Modified:  time.Now().Format(time.RFC3339),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = dat.db.Collection("Reservations").InsertOne(ctx, reservationObj)
	if err != nil {
		dat.logger.Log("Error inserting reservation into DB: ", err.Error())
		return err
	}
	return nil
}
func (dat *database) GetReservation(ctx context.Context, id string) (Reservation, error) {
	tempReservation := Reservation{}
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		dat.logger.Log("Error getting reservation from DB: ", err.Error())
		return tempReservation, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = dat.db.Collection("Reservations").FindOne(ctx, bson.M{"_id": objectID}).Decode(&tempReservation)
	if err != nil {
		dat.logger.Log("Error getting reservation from DB: ", err.Error())
		return tempReservation, err
	}

	return tempReservation, nil
}
func (dat *database) DeleteReservation(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dat.logger.Log("Error deleting reservation from DB: ", err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"_id": objectID}
	res := dat.db.Collection("Reservations").FindOneAndDelete(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		dat.logger.Log("Error deleting reservation from DB: ", err.Error())
		return err
	}
	return nil
}
func (dat *database) GetReservations(ctx context.Context) ([]Reservation, error) {
	tempReservation := Reservation{}
	tempReservations := []Reservation{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := dat.db.Collection("Reservations").Find(ctx, bson.D{})
	if err != nil {
		dat.logger.Log("Error getting reservations from DB: ", err.Error())
		return tempReservations, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err := cursor.Decode(&tempReservation)
		if err != nil {
			dat.logger.Log("Error getting reservation from DB: ", err.Error())
			return tempReservations, err
		}
		tempReservations = append(tempReservations, tempReservation)
	}
	return tempReservations, nil
}
func (dat *database) UpdateReservation(ctx context.Context, id string, from string, to string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		dat.logger.Log("Error updating reservation: ", err.Error())
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"from":     from,
			"to":       to,
			"modified": time.Now().Format(time.RFC3339),
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = dat.db.Collection("Reservations").UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		dat.logger.Log("Error updating reservation: ", err.Error())
		return err
	}

	return nil
}
func (dat *database) GetReservationsFilter(ctx context.Context, chargerID string, userID string) ([]Reservation, error) {
	var mchargerID primitive.ObjectID
	var muserID primitive.ObjectID
	var err error
	filter := bson.M{}
	tempReservation := Reservation{}
	tempReservations := []Reservation{}
	if chargerID != "" {
		mchargerID, err = primitive.ObjectIDFromHex(chargerID)
		if err != nil {
			dat.logger.Log("Error getting reservations from DB: ", err.Error())
			return tempReservations, err
		}
	}
	if userID != "" {
		muserID, err = primitive.ObjectIDFromHex(userID)
		if err != nil {
			dat.logger.Log("Error getting reservations from DB: ", err.Error())
			return tempReservations, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if userID != "" && chargerID == "" {
		filter = bson.M{"userid": muserID}
	} else if userID == "" && chargerID != "" {
		filter = bson.M{"chargerid": mchargerID}
	} else if userID != "" && chargerID != "" {
		filter = bson.M{"chargerid": mchargerID, "userid": muserID}
	}
	cursor, err := dat.db.Collection("Reservations").Find(ctx, filter)
	if err != nil {
		dat.logger.Log("Error getting reservations from DB: ", err.Error())
		return tempReservations, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err := cursor.Decode(&tempReservation)
		if err != nil {
			dat.logger.Log("Error getting reservations from DB: ", err.Error())
			return tempReservations, err
		}
		tempReservations = append(tempReservations, tempReservation)
	}
	return tempReservations, nil
}
