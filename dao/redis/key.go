// Package redis key 的定义
package redis

const (
	KeyPrefix          = "blueblog:"   // KeyPrefix 前缀
	KeyPostTime        = "post:time"   // KeyPostTime zset: 帖子及发帖时间
	KeyPostScore       = "post:score"  // KeyPostScore zset: 帖子及投票
	KeyPostVotedPrefix = "post:voted:" // KeyPostVotedPrefix zset: 记录用户为帖子投票的数据
	KeyCommunityPrefix = "community:"  // KeyCommunityPrefix set: 记录每个社区下的帖子id
)

// 给 redis key 加上前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
