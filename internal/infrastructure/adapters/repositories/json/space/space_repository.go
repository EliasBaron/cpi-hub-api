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

type spaceData struct {
	Spaces []*entity.SpaceEntity `json:"spaces"`
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

	return nil, apperror.NewNotFound("space", name, "space_repository.go:FindByName")
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

	var spaceData spaceData
	if err := json.Unmarshal(data, &spaceData); err != nil {
		r.spaces = make(map[string]*entity.SpaceEntity)
		return
	}

	r.spaces = make(map[string]*entity.SpaceEntity)
	for _, space := range spaceData.Spaces {
		if space != nil && space.ID != "" {
			r.spaces[space.ID] = space
		}
	}
}

func (r *SpaceRepository) saveToFile() error {
	var spaces []*entity.SpaceEntity
	for _, space := range r.spaces {
		spaces = append(spaces, space)
	}

	spaceData := spaceData{
		Spaces: spaces,
	}

	data, err := json.MarshalIndent(spaceData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal spaces: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
