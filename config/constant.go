package config

const (
	AppName = "article"

	RedisTTLOneHour = 1 * 3600 * 1000 * 1000 * 1000

	RedisKeyArticle = "author:%s:search:%s:limit:%d:page:%d"
)
