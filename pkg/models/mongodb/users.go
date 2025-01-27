package mongodb

import "C"
import (
	"context"
	"errors"
	"github.com/aarsinh/auth-golang/pkg/models"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserModel struct {
	C *mongo.Collection
}

func (u *UserModel) GetUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var users []models.User

	cur, err := u.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cur.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserModel) GetUserByID(userID string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	var user models.User

	err = u.C.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserModel) CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return nil, err
	}

	emailExists, checkErr := u.CheckEmailExists(user.Email)
	if checkErr != nil {
		return nil, checkErr
	}

	if emailExists {
		return nil, errors.New("email already in use")
	}

	hashed, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return nil, hashErr
	}
	user.Password = string(hashed)

	// Setting metadata fields
	user.ID = primitive.NewObjectID()
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UserID = user.ID.Hex()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	res, createErr := u.C.InsertOne(ctx, user)
	if createErr != nil {
		return nil, createErr
	}

	return res, nil
}

func (u *UserModel) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var user models.User

	err := u.C.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserModel) CheckEmailExists(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	count, err := u.C.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
