package database

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"log"
	"testjunior/models"
)

var collection *mongo.Collection = Connect(Client, "user")

// Check_session
//1 - если такая сессия есть
//0 - если такой сессии нет
func Check_session(user_Id string) (r int, error error) {
	//Создаем фильтр для поиска по бд, если сессия с таким GUID уже существует
	filter := bson.D{{"user_id", user_Id}}

	var result models.User

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		return 0, error
	}

	if &result != nil {
		fmt.Printf("Found a single document: %+v\n", result)
		return 1, nil
	}
	return 0, nil
}

//
func Save_tokens(hashedToken string, user_Id string) {

	//Создаем модель сессии для бд
	u1 := &models.User{
		Id:            primitive.NewObjectID(),
		Refresh_token: &hashedToken,
		User_id:       user_Id,
	}

	insertResult, err := collection.InsertOne(context.TODO(), u1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}
