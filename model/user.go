package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	Male   = 1
	FeMale = 0
)

type User struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	IdentityID  string        `bson:"identityID" json:"identityID"`
	FirstName   string        `bson:"firstName" json:"firstName"`
	LastName    string        `bson:"lastName" json:"lastName"`
	Birthday    int           `bson:"birthday" json:"birthday"`
	Gender      int           `bson:"gender" json:"gender"`
	PhoneNumber string        `bson:"phoneNumber" json:"phoneNumber"`
	Email       string        `bson:"email" json:"email"`
	Locations   []Location    `bson:"locations" json:"locations"`
	CreatedTime time.Time     `bson:"createdTime" json:"createdTime"`
	UpdatedTime time.Time     `bson:"updatedTime" json:"updatedTime"`
}

func IsValidGender(genderInt int) bool {
	if genderInt != Male && genderInt != FeMale {
		return false
	}
	return true
}
