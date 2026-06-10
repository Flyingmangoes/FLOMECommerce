package redis

import (
	"backend/src/config"
	"backend/src/repository"
	Logger "backend/src/utils/logger"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	glide "github.com/valkey-io/valkey-glide/go/v2"
	glide_config "github.com/valkey-io/valkey-glide/go/v2/config"
	"github.com/valkey-io/valkey-glide/go/v2/options"
	"go.uber.org/zap"
)

type RedisInterface interface {
	GenerateCacheKey(url *url.URL) string
	Set(ctx context.Context, cacheKey, cacheValue string) error
	Get(ctx context.Context, cacheKey string) ([]byte, error)
	Del(ctx context.Context, cacheKey string) error
	Clear(ctx context.Context, cacheKeys []string) error
}

type RedisManager struct {
	Users 			repository.UserStoreInterface
	Products 		repository.ProductStoreInterface
	Stores 			repository.StoreStoreInterface
	REDIS_PORT		string
	REDIS_HOST  	string
	CACHE_TTL		time.Duration
	VClient			*glide.Client
}	

func NewRedisService(
	u repository.UserStoreInterface, 
	p repository.ProductStoreInterface, 
	s repository.StoreStoreInterface, 
	cfg *config.ConfigManager,
) *RedisManager {
	port, err := strconv.Atoi(cfg.REDIS_CONF.REDIS_PORT)
	if err != nil {
		Logger.Log.Error("Failed to get port", zap.Error(err))
		return nil
	}

	vg_cfg := glide_config.NewClientConfiguration().WithAddress(&glide_config.NodeAddress{
		Host: cfg.REDIS_CONF.REDIS_HOST,
		Port: port,
	})

	valkey_client, err := glide.NewClient(vg_cfg)
	if err != nil {
		Logger.Log.Error("Failed to initialize valkey service", zap.Error(err))
		return nil
	}

	return &RedisManager{
		REDIS_PORT: cfg.REDIS_CONF.REDIS_PORT,
		REDIS_HOST: cfg.REDIS_CONF.REDIS_HOST,
		CACHE_TTL: time.Duration(cfg.REDIS_CONF.TTL_M * int(time.Minute)), 
		VClient: valkey_client,
		Users: u,
		Products: p,
		Stores: s,
	}
}

func (rm *RedisManager) ConnectionCheck() {
	res, err := rm.VClient.Ping(context.Background())
	if err != nil {
		Logger.Log.Error("Error", zap.Error(err))
		return
	}

	Logger.Log.Info("Connected", zap.String("response", res))
}

func(rm *RedisManager)Set(ctx context.Context, cacheKey, cacheValue string) error {
	opts := new(options.SetOptions).SetExpiry(&options.Expiry{
		Type: "EX", // in seconds
		Duration: uint64(rm.CACHE_TTL * time.Minute / time.Second),
	})

	res, err := rm.VClient.SetWithOptions(ctx, cacheKey, cacheValue, *opts)
	if err != nil || res.Value() != "OK"{
		return err
	}

	return nil
}

func(rm *RedisManager)Get(ctx context.Context, cacheKey string) ([]byte, error) {
	result, err := rm.VClient.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	return []byte(result.Value()), nil
}

func(rm *RedisManager)Del(ctx context.Context, cacheKey string) error {
	_, err := rm.VClient.Del(ctx , []string{cacheKey})
	return err
}

func(rm *RedisManager)Clear(ctx context.Context, cacheKeys []string) error {
	_, err := rm.VClient.Del(ctx, cacheKeys)
	return err
}

func(rm *RedisManager)GenerateCacheKey(url *url.URL) string {
	return fmt.Sprintf("%s", url)
}