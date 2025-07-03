package config

const (
	AppName = "article"

	RedisErrKeyDoesNotExist = "Key does not exist in Redis"

	RedisTTLOneHour = 1 * 3600 * 1000 * 1000 * 1000

	RedisKeyArticle = "author:%s:search:%s:limit:%d:page:%d"
)
