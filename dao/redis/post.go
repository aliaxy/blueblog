package redis

import (
	"errors"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票的分数
)

var ErrVoteTimeExpire = errors.New("投票时间已过")

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
