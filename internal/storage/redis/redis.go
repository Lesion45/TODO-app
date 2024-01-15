package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

const defExpTime = 0

type Storage struct {
	cl *redis.Client
}

func NewRedisClient() *Storage {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &Storage{cl: redisClient}
}

func (s *Storage) IsRunning() error {
	err := s.cl.Ping()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	return nil
}

// TODO: REVIEW ERRORS
func (s *Storage) AddTask(uuid string, task string) error {
	const op = "storage.redis.AddTask"

	err := s.cl.Set(uuid, task, defExpTime)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}
	return nil
}

// TODO: REVIEW ERRORS
func (s *Storage) GetTask(uuid string) (string, error) {
	fmt.Println(s.cl.Keys("*").Result())
	const op = "storage.redis.GetTask"

	val, err := s.cl.Get(uuid).Result()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return val, nil
}

// TODO: REVIEW ERRORS
func (s *Storage) DeleteTask(uuid string) error {
	const op = "storage.redis.DeleteTask"

	err := s.cl.Del(uuid)
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	return nil
}
