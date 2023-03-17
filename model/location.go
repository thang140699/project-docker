package model

type Location struct {
	Longitude float64 `bson:"longitude" json:"longitude"`
	Latitude  float64 `bson:"latitude" json:"latitude"`
	Address   string  `bson:"address" json:"address"`
	TimeStamp int64   `bson:"timestamp" json:"timestamp"`
	DateStamp int64   `bson:"datestamp" json:"datestamp"`
}
