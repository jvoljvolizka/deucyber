package deucyber

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Config struct {
	DBtype     string `json:"DBtype"`
	DBname     string `json:"DBname"`
	Connection string `json:"Connection"`
}

type newsItem struct {
	Title string `json:"Title"`
	Desc  string `json:"Desc"`
	Link  string `json:"Link"`
	//	Desc string `json:"Desc"`

}

var DBconfig Config
var c *mongo.Client

func GetConf(conf Config) {
	DBconfig = conf
	c = GetClient(DBconfig.Connection)
}

func GetClient(server string) *mongo.Client {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("GetClient recovery error : ' %v ' \n", rec)
		}

	}()

	clientOptions := options.Client().ApplyURI("mongodb://" + server)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	return client
}

func InsertNews(item newsItem) interface{} {
	collection := c.Database(DBconfig.DBname).Collection("News")
	insertResult, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		log.Fatalln("Error on inserting new URL", err)
	}
	return insertResult.InsertedID
}

func GetOneNews(filter bson.M) newsItem {
	var item newsItem
	collection := c.Database(DBconfig.DBname).Collection("News")
	documentReturned := collection.FindOne(context.TODO(), filter)
	documentReturned.Decode(&item)
	return item
}

func GetNews() []*newsItem {
	filter := bson.M{}
	var newsList []*newsItem
	collection := c.Database(DBconfig.DBname).Collection("News")
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
	}
	for cur.Next(context.TODO()) {
		var item newsItem
		err = cur.Decode(&item)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		newsList = append(newsList, &item)
	}
	return newsList
}
