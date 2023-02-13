package testkit

import (
	"context"
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/postgres"

	"gorm.io/gorm/logger"

	"github.com/thunur/gorm-cache/cache"

	"github.com/redis/go-redis/v9"
	"github.com/thunur/gorm-cache/config"
	"gorm.io/gorm"
)

var (
	//postgres://postgres:postgres@192.168.31.230:5432/db_info?sslmode=disable
	username     = "postgres"
	password     = "postgres"
	databaseName = "db_info"
	ip           = "192.168.31.230"
	port         = "5432"
)

var (
	redisIp   = "192.168.31.230"
	redisPort = "6379"
)

var (
	searchCache  *cache.Gorm2Cache
	primaryCache *cache.Gorm2Cache
	allCache     *cache.Gorm2Cache

	searchDB   *gorm.DB
	primaryDB  *gorm.DB
	allDB      *gorm.DB
	originalDB *gorm.DB
)

var (
	testSize = 200 // minimum 200
)

func TestMain(m *testing.M) {
	log("test setup ...")

	var err error
	//logger.Default.LogMode(logger.Info)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, ip, port, databaseName)
	originalDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          logger.Default,
	})
	if err != nil {
		log("open db error: %v", err)
		os.Exit(-1)
	}

	searchDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          logger.Default,
	})

	if err != nil {
		log("open db error: %v", err)
		os.Exit(-1)
	}

	primaryDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          logger.Default,
	})
	if err != nil {
		log("open db error: %v", err)
		os.Exit(-1)
	}

	allDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          logger.Default,
	})
	if err != nil {
		log("open db error: %v", err)
		os.Exit(-1)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: redisIp + ":" + redisPort, DB: 1})

	searchCache, err = cache.NewGorm2Cache(&config.CacheConfig{
		CacheLevel:           config.CacheLevelOnlySearch,
		CacheStorage:         config.CacheStorageRedis,
		RedisConfig:          cache.NewRedisConfigWithClient(redisClient),
		InvalidateWhenUpdate: true,
		CacheTTL:             5000,
		CacheMaxItemCnt:      5000,
		CacheSize:            1000,
		DebugMode:            false,
		Ctx:                  context.Background(),
	})
	if err != nil {
		log("setup search cache error: %v", err)
		os.Exit(-1)
	}

	primaryCache, err = cache.NewGorm2Cache(&config.CacheConfig{
		CacheLevel:           config.CacheLevelOnlyPrimary,
		CacheStorage:         config.CacheStorageRedis,
		RedisConfig:          cache.NewRedisConfigWithClient(redisClient),
		InvalidateWhenUpdate: true,
		CacheTTL:             5000,
		CacheMaxItemCnt:      5000,
		CacheSize:            1000,
		DebugMode:            false,
		Ctx:                  context.Background(),
	})
	if err != nil {
		log("setup primary cache error: %v", err)
		os.Exit(-1)
	}

	allCache, err = cache.NewGorm2Cache(&config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         config.CacheStorageRedis,
		RedisConfig:          cache.NewRedisConfigWithClient(redisClient),
		InvalidateWhenUpdate: true,
		CacheTTL:             5000,
		CacheMaxItemCnt:      5000,
		CacheSize:            1000,
		DebugMode:            true,
		Ctx:                  context.Background(),
	})
	if err != nil {
		log("setup all cache error: %v", err)
		os.Exit(-1)
	}

	primaryDB.Use(primaryCache)
	searchDB.Use(searchCache)
	allDB.Use(allCache)
	// primaryCache.AttachToDB(primaryDB)
	// searchCache.AttachToDB(searchDB)
	// allCache.AttachToDB(allDB)

	err = timer("prepare table and data", func() error {
		return PrepareTableAndData(originalDB)
	})
	if err != nil {
		log("setup table and data error: %v", err)
		os.Exit(-1)
	}

	result := m.Run()

	err = timer("clean table and data", func() error {
		return CleanTable(originalDB)
	})
	if err != nil {
		log("clean table and data error: %v", err)
		os.Exit(-1)
	}

	log("integration test end.")
	os.Exit(result)
}
