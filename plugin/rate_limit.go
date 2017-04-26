package plugin

import (
	"github.com/hellofresh/janus/pkg/api"
	"github.com/hellofresh/janus/pkg/errors"
	"github.com/hellofresh/janus/pkg/middleware"
	"github.com/hellofresh/janus/pkg/router"
	"github.com/hellofresh/janus/pkg/store"
	"github.com/ulule/limiter"
)

const (
	// DefaultPrefix is the default prefix to use for the key in the store.
	DefaultPrefix = "limiter"
)

// RateLimit represents the rate limit plugin
type RateLimit struct {
	storage store.Store
}

// NewRateLimit creates a new instance of HostMatcher
func NewRateLimit(storage store.Store) *RateLimit {
	return &RateLimit{storage}
}

// GetName retrieves the plugin's name
func (h *RateLimit) GetName() string {
	return "rate_limit"
}

// GetMiddlewares retrieves the plugin's middlewares
func (h *RateLimit) GetMiddlewares(config api.Config, referenceSpec *api.Spec) ([]router.Constructor, error) {
	limit := config["limit"].(string)
	policy := config["policy"].(string)

	rate, err := limiter.NewRateFromFormatted(limit)
	if err != nil {
		return nil, err
	}

	limiterStore, err := h.getLimiterStore(policy, referenceSpec.Name)
	if err != nil {
		return nil, err
	}

	return []router.Constructor{
		limiter.NewHTTPMiddleware(limiter.NewLimiter(limiterStore, rate)).Handler,
		middleware.NewRateLimitLogger().Handler,
	}, nil
}

func (h *RateLimit) getLimiterStore(policy string, prefix string) (limiter.Store, error) {
	if prefix == "" {
		prefix = DefaultPrefix
	}

	switch policy {
	case "redis":
		redisStorage, ok := h.storage.(*store.RedisStore)
		if !ok {
			return nil, errors.ErrInvalidStorage
		}

		return limiter.NewRedisStoreWithOptions(redisStorage.Pool, limiter.StoreOptions{
			Prefix:   prefix,
			MaxRetry: limiter.DefaultMaxRetry,
		})
	case "local":
		return limiter.NewMemoryStore(), nil
	default:
		return nil, errors.ErrInvalidPolicy
	}
}
