// Package redis 帖子相关
package redis

import (
	"errors"
	"math"
	"strconv"
	"time"

	"blueblog/models"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600 // 一周的时间 单位秒
	scorePerVote     = 432           // 每一票的分数
)

var (
	// ErrVoteTimeExpire 投票时间已过
	ErrVoteTimeExpire = errors.New("投票时间已过")
	// ErrVoteRepeated 不允许重复投票
	ErrVoteRepeated = errors.New("不允许重复投票")
)

// CreatePost 创建帖子
func CreatePost(postID, communityID int64) error {
	pipeline := client.TxPipeline()
	// 帖子时间
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTime), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子分数
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScore), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 把帖子 ID 加到社区的 set
	cKey := getRedisKey(KeyCommunityPrefix + strconv.FormatInt(communityID, 10))
	pipeline.SAdd(ctx, cKey, postID)

	_, err := pipeline.Exec(ctx)

	return err
}

// VoteForPost 为帖子投票
func VoteForPost(userID, postID string, value float64) error {
	// 判断投票限制
	postTIme := client.ZScore(ctx, getRedisKey(KeyPostTime), postID).Val()
	zap.L().Debug("postTIme", zap.Float64("postTIme", postTIme))
	if float64(time.Now().Unix())-postTIme > oneWeekInSeconds {
		return ErrVoteTimeExpire
	}

	// 更新帖子分数
	// 查投票记录
	oValue := client.ZScore(ctx, getRedisKey(KeyPostVotedPrefix+postID), userID).Val()
	// 不允许重复投票
	if oValue == value {
		return ErrVoteRepeated
	}
	var dir float64
	if value > oValue {
		dir = 1
	} else {
		dir = -1
	}

	diff := math.Abs(oValue - value)

	pipeline := client.TxPipeline()

	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScore), dir*diff*scorePerVote, postID)

	// 记录数据
	if value == 0 {
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedPrefix+postID), userID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedPrefix+postID), redis.Z{
			Score:  value,
			Member: userID,
		})
	}

	_, err := pipeline.Exec(ctx)
	return err
}

// 查询 id
func getIDsFromKey(key string, page, size int64) ([]string, error) {
	// 确定查询的索引起始点
	start := (page - 1) * size
	end := start + size - 1

	// ZREVRANGE 按指定元素从大到小查询指定数量元素
	return client.ZRevRange(ctx, key, start, end).Result()
}

// GetPostIDsInOrder 顺序得到帖子 id
func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从 redis 获取 ID
	key := getRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScore)
	}

	return getIDsFromKey(key, p.Page, p.Size)
}

// GetPostVoteData 根据 ids 查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	// data = make([]int64, 0, len(ids))
	// for _, id := range ids {
	// 	key := getRedisKey(KeyPostVotedPrefix + id)
	// 	// 查找 key
	// 	v := client.ZCount(ctx, key, "1", "1").Val()
	// 	data = append(data, v)
	// }

	// 使用 pipeline 一次发送多条命令 减少 RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedPrefix + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

// GetCommunityPostIDsInOrder 根据社区查询 IDs
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	orderKey := getRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScore)
	}

	// 使用 zinterstore 把分区的帖子 set 与帖子分数的 zset 生成一个新的 zset
	// 对新的 zset 操作

	// 社区的 key
	cKey := getRedisKey(KeyCommunityPrefix + strconv.FormatInt(p.CommunityID, 10))

	key := orderKey + strconv.FormatInt(p.CommunityID, 10)
	if client.Exists(ctx, orderKey).Val() < 1 {
		// 因为是第一次查询，需要根据 post 的 create_time 给 post 分数
		pipeline := client.Pipeline()
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Keys:      []string{cKey, orderKey},
			Aggregate: "MAX",
		})
		pipeline.Expire(ctx, key, 60*time.Second)
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	return getIDsFromKey(key, p.Page, p.Size)
}
