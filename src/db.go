package deucyber

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var DBconfig Config
var c *mongo.Client

func GetConf(conf Config) {
	DBconfig = conf
	c = GetClient(DBconfig.Connection)
}

func DBcheck() {
	for {
		err := c.Ping(context.TODO(), nil)
		if err != nil {
			DBstatus = false
		} else {
			DBstatus = true
		}
		time.Sleep(15000 * time.Millisecond)
	}

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

func InsertNews(item NewsItem) interface{} {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("InsertNews recovery error : ' %v ' \n", rec)
		}

	}()

	collection := c.Database(DBconfig.DBname).Collection("News")
	insertResult, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		log.Panic("Error on inserting new URL", err)
	}
	return insertResult.InsertedID
}

func GetOneNews(filter bson.M) NewsItem {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("GetOneNews recovery error : ' %v ' \n", rec)
		}

	}()

	var item NewsItem
	collection := c.Database(DBconfig.DBname).Collection("News")
	documentReturned := collection.FindOne(context.TODO(), filter)
	documentReturned.Decode(&item)
	return item
}

func GetNews() []*NewsItem {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("GetNews recovery error : ' %v ' \n", rec)
		}

	}()

	filter := bson.M{}
	var newsList []*NewsItem
	collection := c.Database(DBconfig.DBname).Collection("News")
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Panic("Error on Finding all the documents", err)
	}
	for cur.Next(context.TODO()) {
		var item NewsItem
		err = cur.Decode(&item)
		if err != nil {
			log.Panic("Error on Decoding the document", err)
		}
		newsList = append(newsList, &item)
	}
	return newsList
}

func InsertEvents(item EventItem) interface{} {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("InsertEvents recovery error : ' %v ' \n", rec)
		}

	}()

	collection := c.Database(DBconfig.DBname).Collection("Events")
	insertResult, err := collection.InsertOne(context.TODO(), item)
	if err != nil {
		log.Panic("Error on inserting new Event", err)
	}
	return insertResult.InsertedID
}

func GetOneEvent(filter bson.M) EventItem {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("GetOneEvents recovery error : ' %v ' \n", rec)
		}

	}()

	var item EventItem
	collection := c.Database(DBconfig.DBname).Collection("Events")
	documentReturned := collection.FindOne(context.TODO(), filter)
	documentReturned.Decode(&item)
	return item
}

func GetEvents() []*EventItem {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("GetEvents recovery error : ' %v ' \n", rec)
		}

	}()

	filter := bson.M{}
	var eventsList []*EventItem
	collection := c.Database(DBconfig.DBname).Collection("Events")
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Panic("Error on Finding all the documents", err)
	}
	for cur.Next(context.TODO()) {
		var item EventItem
		err = cur.Decode(&item)
		if err != nil {
			log.Panic("Error on Decoding the document", err)
		}
		eventsList = append(eventsList, &item)
	}
	return eventsList
}

func AddAdmin(item Admin) interface{} {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("AddAdmin recovery error : ' %v ' \n", rec)
		}

	}()

	collection := c.Database(DBconfig.DBname).Collection("Admins")
	insertResult, err := collection.InsertOne(context.TODO(), item)
	fmt.Println(insertResult.InsertedID)
	if err != nil {
		log.Panic("Error on inserting new URL", err)
	}
	return insertResult.InsertedID
}

func GetAdmin(filter bson.M) Admin {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("GetAdmin recovery error : ' %v ' \n", rec)
		}

	}()

	var item Admin
	collection := c.Database(DBconfig.DBname).Collection("Admins")
	documentReturned := collection.FindOne(context.TODO(), filter)
	documentReturned.Decode(&item)
	return item
}

func DelAdmin(item Admin) {

	defer func() {
		rec := recover()
		if rec != nil {
			fmt.Printf("DelAdmin recovery error : ' %v ' \n", rec)
		}

	}()

	collection := c.Database(DBconfig.DBname).Collection("Admins")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"id": item.ID})

	if err != nil {
		log.Panic("Error on Deleting admin", err)
	}

}
