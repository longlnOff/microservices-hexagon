package port

import "go.uber.org/zap"

//go:generate mockgen -source=cache.go -destination=mock/cache.go -package=mock

// CacheRepository is an interface for interacting with cache-related business logic
type LoggerRepository interface {
	*zap.Logger
}
