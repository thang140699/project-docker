package main

type Config struct {
	MongoURL string
	RedisURL string
	JwtKey   string
	Binding  string

	
	configPrefix string
	mode         string
	configSource string
}
