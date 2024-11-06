package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

type Models struct {
	LogEntity LogEntity
}

func New(mongoClient *mongo.Client) Models {
	client = mongoClient
	return Models{
		LogEntity: LogEntity{},
	}
}

type LogEntity struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntity) Insert(entry LogEntity) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntity{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error inserting into logs", err)
		return err
	}

	return nil
}

func (l *LogEntity) All() ([]*LogEntity, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()

	opts.SetSort(bson.D{
		bson.E{Key: "created_at", Value: -1},
	})

	cursor, err := collection.Find(ctx, opts)

	if err != nil {
		log.Println("Finding all docs error", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []*LogEntity

	for cursor.Next(ctx) {
		var entity LogEntity

		if err := cursor.Decode(&entity); err != nil {
			log.Println("err decoding log into slice", err)
			return nil, err
		}
		logs = append(logs, &entity)
	}

	return logs, nil
}
