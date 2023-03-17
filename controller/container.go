package main

import (
	database "WeddingBackEnd/database/mongo"
	dbRedis "WeddingBackEnd/database/redis"
	"WeddingBackEnd/database/repository"
	"WeddingBackEnd/ultilities/provider/jwt"
	"WeddingBackEnd/ultilities/provider/mongo"
	"WeddingBackEnd/ultilities/provider/redis"
)

type Provider struct {
	*redis.RedisProvider
	*mongo.MongoProvider
	*jwt.JWTService
}

type Container struct {
	*Provider

	Config                 Config
	UserRepository         repository.UserRepository
	UserAccessIDRepository repository.UserAccessIDRepository
	AccountRepository      repository.AccountRepository
}

func NewContainer(config Config) (*Container, error) {
	container := new(Container)
	err := container.InitContainer(config)
	if err != nil {
		return nil, err
	}

	container.Config = config

	return container, nil
}

func (container *Container) InitContainer(config Config) error {
	// Load providers into container
	err := container.LoadProviders(config)
	if err != nil {
		return err
	}

	// Load repositories
	container.LoadRepositoryImplementations(config)

	return nil
}

func (container *Container) LoadProviders(config Config) error {

	redisProvider := redis.NewRedisProviderFromURL(config.RedisURL)
	mongoProvider := mongo.NewMongoProviderFromURL(config.MongoURL)

	container.Provider = &Provider{
		MongoProvider: mongoProvider,
		RedisProvider: redisProvider,
		JWTService:    jwt.NewJWT(config.JwtKey),
	}
	return nil
}

func (container *Container) LoadRepositoryImplementations(config Config) {
	container.UserRepository = database.NewUserMongoRepository(container.MongoProvider)
	container.AccountRepository = database.NewAccountMongoRepository(container.MongoProvider)
	container.UserAccessIDRepository = dbRedis.NewUserAccessIDRedisRepository(container.RedisProvider)

}
