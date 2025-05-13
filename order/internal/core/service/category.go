package service

import (
	"context"

)

/*
 * CategoryService implements port.CategoryService interface
 * and provides an access to the category repository
 * and cache service
 */
 type CategoryService struct {
	repo  port.CategoryRepository
	cache port.CacheRepository
}

// NewCategoryService creates a new category service instance
func NewCategoryService(repo port.CategoryRepository, cache port.CacheRepository) *CategoryService {
	return &CategoryService{
		repo,
		cache,
	}
}

func (cs *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	category, err := cs.repo.CreateCategory(ctx, category)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	cacheKey := utils.GenerateCacheKey("category", category.ID)
	categorySerialized, err := utils.Serialize(category)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = cs.cache.Set(ctx, cacheKey, categorySerialized, 0)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = cs.cache.DeleteByPrefix(ctx, "categories:*")
	if err != nil {
		return nil, domain.ErrInternal
	}

	return category, nil
}
