package repository

import (
	"context"
	"taskmanagementapi/pkg/config"
	interfaces "taskmanagementapi/pkg/repository/interface"
	"taskmanagementapi/pkg/utils/models"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	UserCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) interfaces.UserRepository {
	return &UserRepository{UserCollection: db.Collection("users")}
}

func (ur *UserRepository) CheckUserExistsByEmail(email string) (bool, error) {
	filter := bson.M{"email": email}
	var result bson.M
	err := ur.UserCollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (ur *UserRepository) UserSignUp(user models.UserSignup) error {
	newUser := bson.M{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}
	_, err := ur.UserCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) FindUserDetailsByEmail(email string) (models.UserDetails, error) {
	filter := bson.M{"email": email}
	var userDetails models.UserDetails
	err := ur.UserCollection.FindOne(context.TODO(), filter).Decode(&userDetails)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.UserDetails{}, nil
		}
		return models.UserDetails{}, err
	}
	return userDetails, nil
}

//GenerateJWTToken
type AuthUserClaims struct {
	Id    string
	Email string
	jwt.StandardClaims
}

func (ur *UserRepository) GenerateJwtToken(user models.UserDetails) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	tokenString, err := GenerateTokenUser(user.ID, user.Email, expirationTime)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateTokenUser(userID string, email string, expirationTime time.Time) (string, error) {
	cfg, _ := config.LoadConfig()
	claims := &AuthUserClaims{
		Id:    userID,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JwtSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

