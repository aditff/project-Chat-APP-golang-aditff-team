// service/user_service.go
package service

import (
	"context"
	"time"

	"project-Chat-APP-golang-aditff-user-service/model"
	"project-Chat-APP-golang-aditff-user-service/repository"

	"github.com/go-redis/redis/v8"
)

type UserService struct {
	Repo  *repository.UserRepository
	Redis *redis.Client
}

const (
	redisOnlineKey = "user:online:" // user:online:<id> -> "1"/"0"
	redisPresenceCh = "presence"    // Pub/Sub channel
)

func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.Repo.GetAll(ctx)
}

func (s *UserService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *UserService) UpdateStatus(ctx context.Context, id string, online bool) error {
	now := time.Now().UTC()
	if err := s.Repo.UpdateStatus(ctx, id, online, now); err != nil { return err }

	// Cache presence sederhana
	val := "0"
	if online { val = "1" }
	if err := s.Redis.Set(ctx, redisOnlineKey+id, val, 0).Err(); err != nil { return err }

	// Publish presence event
	_ = s.Redis.Publish(ctx, redisPresenceCh, id+","+val+","+now.Format(time.RFC3339)).Err()
	return nil
}
