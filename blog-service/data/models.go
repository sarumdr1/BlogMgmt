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
		BlogData: BlogData{},
	}
}

type Models struct {
	BlogData BlogData
}

type BlogData struct {
	ID          string    `json:"id,omitEmpty" bson:"_id,omitempty"`
	Name        string    `json:"name" bson:"name"`
	Author      string    `json:"author" bson:"author"`
	Description string    `json:"description" bson:"description"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

func (t *BlogData) Insert(entry BlogData) (*mongo.InsertOneResult, error) {
	collection := client.Database("blog").Collection("blog")

	blog, err := collection.InsertOne(context.TODO(), BlogData{
		ID:          entry.ID,
		Name:        entry.Name,
		Author:      entry.Author,
		Description: entry.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})

	if err != nil {
		log.Println("Error inserting Blog:", err)
		return nil, err
	}

	return blog, err
}

func (t *BlogData) All() ([]*BlogData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("Blog").Collection("Blog")

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Error finding blog:", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var Blogs []*BlogData

	for cursor.Next(ctx) {
		var item BlogData

		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error decoding Blogs:", err)
			return nil, err
		} else {
			Blogs = append(Blogs, &item)
		}
	}

	return Blogs, nil
}

func (t *BlogData) GetOne(id string) (*BlogData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("blog").Collection("blog")

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting id to ObjectID:", err)
		return nil, err
	}

	var blog BlogData
	err = collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&blog)
	if err != nil {
		log.Println("Error finding Blog:", err)
		return nil, err
	}

	return &blog, nil
}

func (t *BlogData) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("todo").Collection("todo")

	blogId, err := primitive.ObjectIDFromHex(t.ID)
	if err != nil {
		log.Println("Error converting id to ObjectID:", err)
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": blogId},
		bson.D{
			{
				"$set", bson.D{
					{"name", t.Name},
					{"author", t.Author},
					{"description", t.Description},
					{"updated_at", time.Now()},
				}},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (t *BlogData) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("todo").Collection("todo")

	err := collection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping Todo collection:", err)
		return err
	}

	return nil
}

func (t *BlogData) Delete(blog BlogData) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("todo").Collection("todo")

	blogId, err := primitive.ObjectIDFromHex(blog.ID)
	if err != nil {
		log.Println("Error converting id to ObjectID:", err)
		return err
	}

	_, err = collection.DeleteOne(ctx, bson.M{"_id": blogId})
	if err != nil {
		log.Println("Error deleting blog:", err)
		return err
	}

	return nil
}
