package database

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"log"
	"testjunior/models"
)

var collection *mongo.Collection = Connect(Client, "user")

// Check_session
//1 - если есть пользователь с таким GUID
//0 - если пользователя с таким GUID нет
func Check_session(user_Id string) (r int, error error, result models.User) {
	//Создаем фильтр для поиска по бд, если сессия с таким GUID уже существует
	filter := bson.D{{"user_id", user_Id}}

	//var result models.User

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		return 0, error, result
	}

	if &result != nil {
		fmt.Printf("Found a single document: %+v\n", result)
		return 1, nil, result
	}
	return 0, nil, result
}

func Check_refresh_token(refresh_token string, result models.User) (r bool, error error) {
	if result.Refresh_token != nil {
		err := bcrypt.CompareHashAndPassword([]byte(*result.Refresh_token), []byte(refresh_token))
		fmt.Println("**************************")
		fmt.Println(*result.Refresh_token)
		fmt.Println(refresh_token)
		fmt.Println(err)
		fmt.Println("**************************")
		if err != nil {
			fmt.Println("false")
			return false, err
		}
		fmt.Println("true")
		return true, nil
	}
	fmt.Println("3")
	return false, errors.New("empty db User")
}

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

func DeleteClientSessionByGUID(user_id string) {
	filter := bson.D{{"user_id", user_id}}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}