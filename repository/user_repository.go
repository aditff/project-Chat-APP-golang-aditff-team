// repository/user_repository.go
package repository

import (
	"context"
	"time"

	"project-Chat-APP-golang-aditff-user-service/model"

	"gorm.io/gorm"
)

type UserRepository struct{ DB *gorm.DB }

func (r *UserRepository) GetAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	if err := r.DB.WithContext(ctx).Order("name").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	if err := r.DB.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UpdateStatus(ctx context.Context, id string, online bool, at time.Time) error {
	data := map[string]interface{}{
		"online": online,
	}
	if !online {
		data["last_seen"] = at
	}
	return r.DB.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Updates(data).Error
}
