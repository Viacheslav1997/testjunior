package database

import (
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"log"
	"testjunior/models"
)

var collection *mongo.Collection = Connect(Client, DbCollectionName)

func Check_session(user_Id string, refresh_token string) (r bool, error error, result models.User) {
	//Создаем фильтр для поиска по бд, если сессия с таким GUID уже существует
	filter := bson.D{{"user_id", user_Id}}
	options := options.Find()
	options.SetLimit(10)

	cur, err := collection.Find(context.TODO(), filter, options)

	if err != nil {
		return false, error, result
	}

	for cur.Next(context.TODO()) {
		fmt.Println("Checking...")
		var elem models.User
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		//Проверяем все сессии, на наличие refresh токена
		err = bcrypt.CompareHashAndPassword([]byte(*elem.Refresh_token), []byte(refresh_token))

		if err == nil {
			cur.Close(context.TODO())
			fmt.Println("Found document")
			return true, nil, elem
		}
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	fmt.Println("Document not found")

	return false, nil, result
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

func DeleteClientSessionByGUID(user_id string, refresh_token *string) {
	filter := bson.D{{"user_id", user_id}, {"refresh_token", refresh_token}}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
}

func ValidateSession(refreshToken string, accessTokenString string) bool {

	decodedRefreshToken, _ := base64.StdEncoding.DecodeString(refreshToken)
	var refreshTokenUniq interface{}

	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(string(decodedRefreshToken), claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("Uniq"), nil
	})

	for key, val := range claims {
		if key == "Uniq" {
			refreshTokenUniq = val
		}
	}

	var accessTokenUniq interface{}

	claims2 := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(accessTokenString, claims2, func(token *jwt.Token) (interface{}, error) {
		return []byte("Uniq"), nil
	})

	for key, val := range claims2 {
		if key == "Uniq" {
			accessTokenUniq = val
		}
	}

	if refreshTokenUniq == accessTokenUniq {
		return true
	}
	return false

}
