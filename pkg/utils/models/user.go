package models

type UserSignup struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password"  validate:"required,min=6,max=20"`
}

type UserSignIn struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=6,max=20"`
}

type UserDetails struct {
	ID       string `bson:"_id"`
	Name     string `bson:"name"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
}
