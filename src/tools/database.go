package tools

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DB struct {
	client               *mongo.Client
	moderationCollection *mongo.Collection
	configCollection     *mongo.Collection
}

var dbInstance *DB

// Инициализация соединения с базой данных
func InitDB() error {

	MONGO_DB_CONNECTION := os.Getenv("MONGO_DB_CONNECTION")
	clientOptions := options.Client().ApplyURI(MONGO_DB_CONNECTION)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("ошибка подключения к бд: %w", err)
		return fmt.Errorf("ошибка подключения к бд: %w", err)
	}

	DB_NAME := os.Getenv("MONGO_DB_NAME")

	moderationCollection := client.Database(DB_NAME).Collection("moderation")
	configCollection := client.Database(DB_NAME).Collection("config")

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Printf("failed to ping mongo: %s", err.Error())
		return fmt.Errorf("failed to ping mongo: %w", err)
	}

	dbInstance = &DB{
		client:               client,
		moderationCollection: moderationCollection,
		configCollection:     configCollection,
	}

	return nil
}

func Get_Punish(filter bson.M) (*DbPunish, error) {
	var punish *DbPunish
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := dbInstance.moderationCollection.FindOne(ctx, filter).Decode(&punish)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}
	return punish, nil
}

func Get_Config(guildId string) (*DbConfig, error) {
	var config DbConfig
	filter := bson.M{
		"guildId": guildId,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := dbInstance.configCollection.FindOne(ctx, filter).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, fmt.Errorf("failed to decode document: %w", err)
	}
	return &config, nil
}

func Update_Config(guildId string, data bson.M) error {
	filter := bson.M{"guildId": guildId}
	updateData := bson.M{"$set": data}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := dbInstance.configCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return err
	}
	return nil
}

func Update_Punish(filter bson.M, data bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	updateData := bson.M{"$set": data}
	defer cancel()
	_, err := dbInstance.moderationCollection.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return err
	}
	return nil
}

func Insert_Config(data bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := dbInstance.configCollection.InsertOne(ctx, data)
	if err != nil {
		return err
	}
	return nil
}

func Insert_Punish(data bson.M) error {
	_, err := dbInstance.moderationCollection.InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}
	return nil
}

func Get_Punishments(filter bson.M) ([]*DbPunish, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cursor, err := dbInstance.moderationCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var punishments []*DbPunish
	if err := cursor.All(ctx, &punishments); err != nil {
		return nil, err
	}

	return punishments, nil
}
