package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Account struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	Scopes    string        `bson:"scopes" json:"scopes"`
	Email     string        `bson:"email" json:"email"`
	Password  string        `bson:"password,omitempty" json:"-"`
	UserID    bson.ObjectId `bson:"userID" json:"userID"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
	CreatedBy string        `bson:""createdBy json:"createBy"`
}
