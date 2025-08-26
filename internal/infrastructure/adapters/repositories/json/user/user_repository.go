package user

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"cpi-hub-api/internal/core/domain"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/json/entity"
	"cpi-hub-api/internal/infrastructure/adapters/repositories/json/mapper"
	"cpi-hub-api/pkg/apperror"
)

type UserRepository struct {
	filePath string
	users    map[string]*entity.UserEntity
	mutex    sync.RWMutex
}

type userData struct {
	Users []*entity.UserEntity `json:"users"`
}

func NewUserRepository(filePath string) *UserRepository {
	repo := &UserRepository{
		filePath: filePath,
		users:    make(map[string]*entity.UserEntity),
	}

	repo.loadFromFile()

	return repo
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	userEntity := mapper.ToJSONDatabaseUser(user)

	r.users[userEntity.ID] = userEntity

	if err := r.saveToFile(); err != nil {
		delete(r.users, user.ID)
		return fmt.Errorf("failed to save user to file: %w", err)
	}

	return nil
}

func (r *UserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.ID == id {
			return mapper.ToDomainUser(user), nil
		}
	}

	return nil, apperror.NewNotFound("User not found", nil, "user_repository.go:FindById")
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user := r.users[email]

	if user == nil {
		return nil, nil
	}

	return mapper.ToDomainUser(user), nil
}

func (r *UserRepository) loadFromFile() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		r.users = make(map[string]*entity.UserEntity)
		return
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		r.users = make(map[string]*entity.UserEntity)
		return
	}

	var userData userData
	if err := json.Unmarshal(data, &userData); err != nil {
		r.users = make(map[string]*entity.UserEntity)
		return
	}

	r.users = make(map[string]*entity.UserEntity)
	for _, user := range userData.Users {
		if user != nil && user.Email != "" {
			r.users[user.Email] = user
		}
	}
}

func (r *UserRepository) saveToFile() error {

	var db entity.DatabaseFile

	if data, err := os.ReadFile(r.filePath); err == nil {
		_ = json.Unmarshal(data, &db)
	}

	db.Users = nil
	for _, user := range r.users {
		db.Users = append(db.Users, user)
	}

	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
