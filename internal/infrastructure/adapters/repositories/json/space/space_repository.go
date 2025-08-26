package space

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

type SpaceRepository struct {
	filePath string
	spaces   map[string]*entity.SpaceEntity
	mutex    sync.RWMutex
}

func NewSpaceRepository(filePath string) *SpaceRepository {
	repo := &SpaceRepository{
		filePath: filePath,
		spaces:   make(map[string]*entity.SpaceEntity),
	}

	repo.loadFromFile()

	return repo
}

func (r *SpaceRepository) Create(ctx context.Context, space *domain.Space) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	spaceEntity := mapper.ToJSONDatabaseSpace(space)
	r.spaces[spaceEntity.ID] = spaceEntity

	if err := r.saveToFile(); err != nil {
		delete(r.spaces, space.ID)
		return fmt.Errorf("failed to save space to file: %w", err)
	}

	return nil
}

func (r *SpaceRepository) FindById(ctx context.Context, id string) (*domain.Space, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, space := range r.spaces {
		if space.ID == id {
			return mapper.ToDomainSpace(space), nil
		}
	}

	return nil, apperror.NewNotFound("space", id, "space_repository.go:FindById")
}

func (r *SpaceRepository) FindByName(ctx context.Context, name string) (*domain.Space, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, space := range r.spaces {
		if space.Name == name {
			return mapper.ToDomainSpace(space), nil
		}
	}

	return nil, nil
}

func (r *SpaceRepository) loadFromFile() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, err := os.Stat(r.filePath); os.IsNotExist(err) {
		r.spaces = make(map[string]*entity.SpaceEntity)
		return
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		r.spaces = make(map[string]*entity.SpaceEntity)
		return
	}

	var db entity.DatabaseFile
	if err := json.Unmarshal(data, &db); err != nil {
		r.spaces = make(map[string]*entity.SpaceEntity)
		return
	}

	r.spaces = make(map[string]*entity.SpaceEntity)
	for _, space := range db.Spaces {
		if space != nil && space.ID != "" {
			r.spaces[space.ID] = space
		}
	}
}

func (r *SpaceRepository) saveToFile() error {
	// NO uses r.mutex.Lock() aquí

	// Leer el archivo actual para no perder los usuarios
	var db entity.DatabaseFile
	if data, err := os.ReadFile(r.filePath); err == nil {
		_ = json.Unmarshal(data, &db)
	}

	// Actualizar solo los espacios
	db.Spaces = nil
	for _, space := range r.spaces {
		db.Spaces = append(db.Spaces, space)
	}

	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal spaces: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
