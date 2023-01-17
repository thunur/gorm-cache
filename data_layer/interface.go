package data_layer

import (
	"context"

	"github.com/thunur/gorm-cache/config"
	"github.com/thunur/gorm-cache/util"
)

type DataLayerInterface interface {
	Init(conf *config.CacheConfig, prefix string, ctx context.Context) error

	// read
	BatchKeyExist(ctx context.Context, keys []string) (bool, error)
	KeyExists(ctx context.Context, key string) (bool, error)
	GetValue(ctx context.Context, key string) (string, error)
	BatchGetValues(ctx context.Context, keys []string) ([]string, error)

	// write
	CleanCache(ctx context.Context) error
	DeleteKeysWithPrefix(ctx context.Context, keyPrefix string) error
	DeleteKey(ctx context.Context, key string) error
	BatchDeleteKeys(ctx context.Context, keys []string) error
	BatchSetKeys(ctx context.Context, kvs []util.Kv) error
	SetKey(ctx context.Context, kv util.Kv) error
}
