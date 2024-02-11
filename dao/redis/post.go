package redis

import (
	"errors"
	"math"
	"time"

	"blueblog/models"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票的分数
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeated   = errors.New("不允许重复投票")
)

func CreatePost(postID int64) error {
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

	_, err := pipeline.Exec(ctx)

	return err
}

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

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从 redis 获取 ID
	key := getRedisKey(KeyPostTime)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScore)
	}

	// 确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1

	// ZREVRANGE 按指定元素从大到小查询指定数量元素
	return client.ZRevRange(ctx, key, start, end).Result()
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
