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
	err := s.cl.Ping().Err()
	if err != nil {
		return fmt.Errorf("%s", err)
	}
	return nil
}

func (s *Storage) AddTask(uuid string, task string) error {
	const op = "storage.redis.AddTask"

	status, err := s.cl.Exists("users:" + uuid).Result()
	if err != nil {
		return fmt.Errorf("%s: %s", op, err)
	}

	if status == 1 {
		err := s.cl.HSet("users:"+uuid, task, true).Err()
		if err != nil {
			return fmt.Errorf("%s: %s", op, err)
		}
	} else {
		err := s.cl.HMSet("users:"+uuid, map[string]interface{}{uuid: true}).Err()
		if err != nil {
			return fmt.Errorf("%s: %s", op, err)
		}
	}
	return nil
}

func (s *Storage) GetTasks(uuid string) (map[string]string, error) {
	const op = "storage.redis.GetTask"

	val, err := s.cl.HGetAll("users:" + uuid).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	return val, nil
}

func (s *Storage) DeleteTask(uuid string, task string) error {
	const op = "storage.redis.DeleteTask"

	_, err := s.cl.HDel("user:"+uuid, task).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("%s: %w", op, "task does not exist")
		} else {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
