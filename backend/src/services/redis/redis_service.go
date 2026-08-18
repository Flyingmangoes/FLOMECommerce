package cache_service

import (
	"backend/src/config"
	logger_system "backend/src/utils/LoggerSystem"
	"context"
	"fmt"
	"net"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisInterface interface {
	GenerateCacheKey(prefix string, rawQuery string) string
	Set(ctx context.Context, cacheKey string, cacheValue string) error
	Get(ctx context.Context, cacheKey string) ([]byte, error)
	Del(ctx context.Context, cacheKey string) error
	Clear(ctx context.Context, cacheKeys []string) error
}

type RedisManager struct {
	REDIS_PORT		string
	REDIS_HOST  	string
	CACHE_TTL		time.Duration
	Client			*redis.Client
}	

func NewRedisService(cfg *config.ConfigManager) *RedisManager {
	redis_client:= redis.NewClient(&redis.Options{
		Addr: net.JoinHostPort(cfg.REDIS_CONF.REDIS_HOST, cfg.REDIS_CONF.REDIS_PORT),
	})

	return &RedisManager{
		REDIS_PORT: cfg.REDIS_CONF.REDIS_PORT,
		REDIS_HOST: cfg.REDIS_CONF.REDIS_HOST,
		CACHE_TTL: time.Duration(cfg.REDIS_CONF.CACHE_TTL * int(time.Minute)), 
		Client: redis_client,
	}
}

func (rm *RedisManager) ConnectionCheck() {
	err := rm.Client.Ping(context.Background()).Err()
	if err != nil {
		logger_system.Log.Error("Failed to connect to redis", zap.Error(err))
	}

	logger_system.Log.Info("Connected", zap.String("response", "OK"))
}

func(rm *RedisManager)Set(ctx context.Context, cacheKey string, cacheValue string) error {
	return rm.Client.Set(ctx, cacheKey, cacheValue, rm.CACHE_TTL).Err()
}

func(rm *RedisManager)Get(ctx context.Context, cacheKey string) ([]byte, error) {
	val, err := rm.Client.Get(ctx, cacheKey).Bytes()
	return val, err
}

func(rm *RedisManager)Del(ctx context.Context, cacheKey string) error {
	return rm.Client.Del(ctx , cacheKey).Err()
}

func(rm *RedisManager)Clear(ctx context.Context, cacheKeys []string) error {
	return rm.Client.Del(ctx, cacheKeys...).Err()
}

func(rm *RedisManager)GenerateCacheKey(prefix, rawQuery string) string {
	params, _ := url.ParseQuery(rawQuery)

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params.Get(k)))
	}

	return fmt.Sprintf("%s:%s", prefix, strings.Join(parts, ":"))
}
