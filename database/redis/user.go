package redis

import (
	"WeddingBackEnd/database/repository"
	"WeddingBackEnd/ultilities/provider/redis"
	"fmt"
	"strconv"
	"time"
)

var UserRedisAccessIdPrefix = "user_access"

type UserAccessIDRedisRepository struct {
	redisProvider *redis.RedisProvider
	redisClient   *redis.RedisClient
}

func NewUserAccessIDRedisRepository(provider *redis.RedisProvider) repository.UserAccessIDRepository {
	return &UserAccessIDRedisRepository{redisProvider: provider, redisClient: provider.RedisClient()}
}

func (repo *UserAccessIDRedisRepository) redisKey(username string) string {
	return fmt.Sprintf("%s:%s", UserRedisAccessIdPrefix, username)
}

func (repo *UserAccessIDRedisRepository) Set(username, accessID, exp string) error {
	return repo.redisClient.HSet(repo.redisKey(username), accessID, exp)
}

func (repo *UserAccessIDRedisRepository) Exists(username, accessID string) (bool, error) {
	return repo.redisClient.HExists(repo.redisKey(username), accessID)
}

func (repo *UserAccessIDRedisRepository) Get(username, accessID string) (result time.Time, err error) {
	val, err := repo.redisClient.HGet(repo.redisKey(username), accessID)
	if err != nil {
		return
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return
	}

	return time.Unix(int64(i), 0), nil
}

func (repo *UserAccessIDRedisRepository) All(username string) (map[string]string, error) {
	return repo.redisClient.HGetAll(repo.redisKey(username))
}

func (repo *UserAccessIDRedisRepository) Delete(username, accessID string) error {
	return repo.redisClient.HDel(repo.redisKey(username), accessID)
}
