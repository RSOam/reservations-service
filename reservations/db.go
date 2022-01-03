package reservations

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net/http"
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
func (dat *database) ReservationClosest(ctx context.Context, userID string, from string, to string, location Location) (Reservation, error) {
	tempReservation := Reservation{}
	userIDmongo, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		dat.logger.Log("Error creating reservation: ", err.Error())
		return tempReservation, err
	}
	requestBody, _ := json.Marshal(GetChargersRequest{})
	client := &http.Client{}
	chargersAddr, _ := getConsulValue(dat.consul, dat.logger, "chargersService")
	chargersUri := chargersAddr + "/chargers"
	req, err := http.NewRequest(http.MethodGet, chargersUri, bytes.NewBuffer(requestBody))
	if err != nil {
		return tempReservation, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return tempReservation, err
	}
	defer resp.Body.Close()
	tempResponse := GetChargersResponse{}
	err = json.NewDecoder(resp.Body).Decode(&tempResponse)
	if err != nil {
		return tempReservation, err
	}
	client.CloseIdleConnections()
	tmpChar, _ := getClosestCharger(location, tempResponse.Chargers)
	reservationObj := Reservation{
		ChargerID: tmpChar.ID,
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
		return tempReservation, err
	}
	tempReservation.ChargerID = tmpChar.ID
	tempReservation.UserID = userIDmongo
	tempReservation.From = from
	tempReservation.To = to
	tempReservation.Created = reservationObj.Created
	tempReservation.Modified = reservationObj.Modified
	return tempReservation, nil
}
func getClosestCharger(location Location, chargers []Charger) (Charger, error) {
	closest := Charger{}
	minDist := 100000.0
	for _, charger := range chargers {
		dst, _ := calcDistance(location, charger.Location)
		if dst < minDist {
			closest = charger
			minDist = dst
		}
	}
	return closest, nil
}
func getConsulValue(consul consulapi.Client, logger log.Logger, key string) (string, error) {
	kv := consul.KV()
	keyPair, _, err := kv.Get(key, nil)
	if err != nil {
		logger.Log("msg", "Failed getting consul key")
		return "", err
	}
	return string(keyPair.Value), nil
}
func calcDistance(loc1 Location, loc2 Location) (float64, error) {
	radlat1 := float64(math.Pi * loc1.Latitude / 180)
	radlat2 := float64(math.Pi * loc2.Latitude / 180)
	theta := float64(loc1.Longitude - loc2.Longitude)
	radtheta := float64(math.Pi * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515 * 1.60934

	return dist, nil
}
