package models

import (
	"context"
	"fmt"
	"log"
	"sanchit/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	IsActive bool               `json:"isActive" bson:"isActive"`
}

/*user value -- fetching object will be long
user value stored at memory address -- fetching object will be faster
*/

func UserCreate(newUser *User) (*User, error) {
	// Password Hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return &User{}, err
	}
	//Collection := config.MongoBlog.Collection("users")
	Collection := config.GetCollection("user_data")
	//It will create the default table users table with the help of mongoDb when nothing exists in the database.
	newUser.ID = primitive.NewObjectID()
	newUser.Password = string(hashedPassword)
	newUser.IsActive = true
	_, err = Collection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return &User{}, err
	}
	return newUser, nil
}

func UserbyId(id string) (*User, error) {
	var newUser User
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &User{}, err
	}
	Collection := config.GetCollection("user_data")
	err = Collection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectId}}).Decode(&newUser)
	if err != nil {
		return &User{}, err
	}
	return &newUser, nil
}

func UserbyEmail(email string) (*User, error) {
	var newUser User
	Collection := config.GetCollection("user_data")
	err := Collection.FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(&newUser)
	if err != nil {
		return &User{}, err
	}
	return &newUser, nil
}

func GetAllUsers() (*[]User, error) {
	Collection := config.GetCollection("user_data")
	//c.IndentedJSON(http.StatusOK, books)
	cursor, err := Collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return &[]User{}, err
	}
	var users []User
	for cursor.Next(context.TODO()) {
		var output User
		if err := cursor.Decode(&output); err != nil {
			log.Fatal(err)
		}
		//fmt.Println(result)
		users = append(users, output)
	}
	return &users, err
}

// fe_user - front end user
func VerifyUser(fe_user *User) bool {
	db_user, err := UserbyEmail(fe_user.Email)
	if err != nil {
		fmt.Println("User Not Found")
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(db_user.Password), []byte(fe_user.Password))
	if err != nil {
		fmt.Println("User Authentication Failed")
		return false
	} else {
		fmt.Println("User Authenticated Successfully")
		return true
	}
}
