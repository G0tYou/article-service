package redis

import (
	"article/config"
	"article/pkg/listing"
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

type Storage struct {
	db *redis.Client
}

func NewStorage(cfgdb config.Redis) (*Storage, error) {
	var err error

	s := new(Storage)

	opts := &redis.Options{Addr: cfgdb.Addr, DB: cfgdb.Db}
	s.db = redis.NewClient(opts)

	_, err = s.db.Ping(context.Background()).Result()
	if err != nil {
		return s, err
	}

	log.Println("Redis connected", cfgdb.Addr, cfgdb.Db)

	return s, nil
}

func (s *Storage) CreateArticle(ctx context.Context, lfgar listing.FilterGetArticle) {
	cg, _ := json.Marshal(lfgar)

	s.db.Set(ctx, config.RedisKeyArticle, cg, config.RedisTTLOneHour)
}

func (s *Storage) DeleteArticle(ctx context.Context, aid int) {
	s.db.FlushDB(ctx)
}
