package feeds

import (
	"context"
	feeds "github.com/rfomin84/discrep-service/internal/services/feeds/domain"
)

type StoreInterface interface {
	Save(ctx context.Context, key string, data interface{}) error
	Get(ctx context.Context, key string) ([]feeds.Feed, error)
}
