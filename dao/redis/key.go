package redis

const (
	KeyPrefix          = "blueblog:"
	KeyPostTime        = "post:time"   // zset: 帖子及发帖时间
	KeyPostScore       = "post:score"  // zset: 帖子及投票
	KeyPostVotedPrefix = "post:voted:" // zset: 记录用户为帖子投票的数据
)

// getRedisKey 给 redis key 加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
