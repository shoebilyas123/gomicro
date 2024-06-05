package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}


type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID string `bson:"_id,omitempty",json:"id,omitempty"`
	Name string `bson:"name",json:"name"`
	Data string `bson:"data",json:"data"`
	CreatedAt time.Time `bson:"created_at",json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at",json:"updated_at"`
}

func (l *LogEntry) InsertOne() error {
	collection := client.Database("logs").Collection("logs");

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name: l.Name,
		Data: l.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	});

	if err != nil {
		log.Println("Error inserting entry into the DB,",err);
		return err;
	}

	return nil;
}

func (l *LogEntry) GetAll() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second);
	defer cancel()

	collection := client.Database("logs").Collection("logs");
	ops := options.Find();

	ops.SetSort(bson.D{{"created_at",-1}});

	cursor, err := collection.Find(context.TODO(), bson.D{}, ops);

	if err != nil {
		log.Println("Finding all docs error: ",err);
		return nil, err;
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)

		if err != nil {
			log.Println("Error decoding log into slice: ", err);
			return nil, err
		}

		logs = append(logs, &item);
	}


	return logs, nil;
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second);
	defer cancel()

	collection := client.Database("logs").Collection("logs");

	docId, err := primitive.ObjectIDFromHex(id);

	if err != nil {
		log.Println("Error getting ObjectID from Hex: ",err);

		return nil, err;
	}

	var logEntry LogEntry;

	err = collection.FindOne(ctx, bson.M{"_id":docId}).Decode(&logEntry);

	if err != nil {
		log.Println("Error fetching log entry from DB: ", err);
		return nil, err;
	}
	
	return &logEntry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second);
	defer cancel()

	collection := client.Database("logs").Collection("logs");

	if err := collection.Drop(ctx); err != nil {
		log.Println("Error dropping colleciton: ",err);
		return err;
	}

	return nil;
}

func (l *LogEntry) UpdateResult() (*mongo.UpdateResult,error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second);
	defer cancel()

	collection := client.Database("logs").Collection("logs");

	docId, err := primitive.ObjectIDFromHex(l.ID)

	if err != nil {
		log.Println("Error getting ObjectID from Hex: ", err);
		return nil, err
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": docId}, bson.D{
		{
			"$set", bson.D{
				{"name", l.Name},
				{"data", l.Data},
				{"updatedAt", time.Now()},
			},
		},
})

	if err != nil {
		log.Println("Error updating doc: ", err);
		return nil, err
	}

	return result, nil;
}