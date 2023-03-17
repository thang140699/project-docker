package model

type GoogleUserDetails struct {
	ID      string `bson:"id" json:"id"`
	Email   string `bson:"email" json:"email"`
	Picture string `bson:"picture" json:"picture"`
}
