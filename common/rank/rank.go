package rank

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/wuyfueng/rank/common/constants"
	"github.com/wuyfueng/rank/common/redis_wrapper"
)

// Expire 设置过期时间 second: 剩余秒数
func (rc *RankConf) Expire(regionId int64, second int) (err error) {
	err = redis_wrapper.Rdb().Expire(context.TODO(), rc.RedisKey(regionId), time.Second*time.Duration(second)).Err()
	return
}

// Del 删除排行榜
func (rc *RankConf) Del(regionId int64) (err error) {
	err = redis_wrapper.Rdb().Del(context.TODO(), rc.RedisKey(regionId)).Err()
	return
}

// Remove 从排行榜移除某个元素
func (rc *RankConf) Remove(regionId int64, member string) (ret int64, err error) {
	ret, err = redis_wrapper.Rdb().ZRem(context.TODO(), rc.RedisKey(regionId), member).Result()
	return
}

// AddWait 添加到待更新
func (rc *RankConf) AddWait(regionId int64, member string) {
	err := redis_wrapper.Rdb().SAdd(context.TODO(), rc.WaitRedisKey(), fmt.Sprintf("%d_%s", regionId, member)).Err()
	if err != nil {
		// err_log
	}
}

// RemoveWait 从待更新移除
func (rc *RankConf) RemoveWait(member string) {
	err := redis_wrapper.Rdb().SRem(context.TODO(), rc.WaitRedisKey(), member).Err()
	if err != nil {
		// err_log
	}
}

// GetRank 查询排名[1,2,3...]
func (rc *RankConf) GetRank(regionId int64, member string) (rank int64, err error) {
	if rc.IsPositiveSort() { // 是否正序
		rank, err = redis_wrapper.Rdb().ZRank(context.TODO(), rc.RedisKey(regionId), member).Result()
		if err != nil {
			return 0, err
		}
	} else {
		rank, err = redis_wrapper.Rdb().ZRevRank(context.TODO(), rc.RedisKey(regionId), member).Result()
		if err != nil {
			return 0, err
		}
	}

	rank++

	return
}

// GetScore 查询分数
func (rc *RankConf) GetScore(regionId int64, member string) (score int64, err error) {
	fScore, err := redis_wrapper.Rdb().ZScore(context.TODO(), rc.RedisKey(regionId), member).Result()
	if err != nil {
		return
	}

	// 分数是否带时间
	score = int64(fScore)
	if rc.IsTimeScore() {
		score = GetRankRealScore(score)
	}
	return
}

// GetRedisScore 查询redis原始分数
func (rc *RankConf) GetRedisScore(regionId int64, member string) (score int64, err error) {
	fScore, err := redis_wrapper.Rdb().ZScore(context.TODO(), rc.RedisKey(regionId), member).Result()
	if err != nil {
		return
	}

	score = int64(fScore)
	return
}

// GetFirstScore 查询第一名及分数
func (rc *RankConf) GetFirstScore(regionId int64) (member string, score int64, err error) {
	// 获取第一名的分数
	zl, err := redis_wrapper.Rdb().ZRangeWithScores(context.TODO(), rc.RedisKey(regionId), 0, 0).Result()
	if err != nil {
		return "", 0, err
	}
	if len(zl) != 1 {
		return "", 0, errors.New("zl is empty")
	}

	member, score, err = ParseUserScoreRet(zl[0], rc.IsTimeScore())
	if err != nil {
		return "", 0, err
	}

	return
}

// GetRankScore 查询第x名(1,2,3...)及分数
func (rc *RankConf) GetRankScore(regionId int64, rank int64) (member string, score int64, err error) {
	// 获取第一名的分数
	zl, err := redis_wrapper.Rdb().ZRevRangeWithScores(context.TODO(), rc.RedisKey(regionId), rank-1, rank-1).Result()
	if err != nil {
		return "", 0, err
	}
	if len(zl) != 1 {
		return "", 0, errors.New("zl is empty")
	}

	member, score, err = ParseUserScoreRet(zl[0], rc.IsTimeScore())
	if err != nil {
		return "", 0, err
	}

	return
}

// RankListRange 遍历排行榜
// f func(memberId int64, rank int, score int64)
//  rank: 1,2,3....
//  score: 不带时间戳的分数
func (rc *RankConf) RankListRange(regionId int64, top int64, f func(member string, rank int, score float64)) (err error) {
	if top < 1 {
		return errors.New("top < 1")
	}

	zl, err := redis_wrapper.Rdb().ZRevRangeWithScores(context.TODO(), rc.RedisKey(regionId), 0, top-1).Result()
	if err != nil {
		return err
	}

	for k, v := range zl {
		f(v.Member.(string), k, v.Score)
	}

	return
}

// RankListRangeWithRawScore 遍历排行榜
// f func(memberId int64, rank int, score int64)
//  rank: 1,2,3....
//  rawScore: 原生分数
func (rc *RankConf) RankListRangeWithRawScore(regionId int64, top int64, f func(member string, rank int, rawScore float64)) (err error) {
	zl, err := redis_wrapper.Rdb().ZRevRangeWithScores(context.TODO(), rc.RedisKey(regionId), 0, top-1).Result()
	if err != nil {
		return err
	}

	for k, v := range zl {
		f(v.Member.(string), k, v.Score)
	}

	return
}

// ParseUserScoreRet 解析单个排行榜结果 前端展示分数
func ParseUserScoreRet(z redis.Z, isTimeScore bool) (member string, score int64, err error) {
	member = z.Member.(string)
	score = int64(z.Score)

	// 分数是否带时间
	if isTimeScore {
		score = GetRankRealScore(score)
	}

	return
}

// RankTimeScoreZIncrByAtomic 排行榜带时间分数 原子incr
// KEYS[1] redisKey
// ARGV[1] userId member
// ARGV[2] incrScore
// ARGV[3] offsetBits 偏移位数
// ARGV[4] complementTime 分数追加值(时间补值)
// return 当前分数
func RankTimeScoreZIncrByAtomic(redisKey string, member string, incrScore, offsetBits, complementTime int64) (currentScore int64, err error) {
	script := redis.NewScript(`
	local score = 0; 
	local ret = redis.call('ZSCORE', KEYS[1], ARGV[1]);
	if ret then
		score = tonumber(ret);
		local n1, n2 = math.modf(score/(2^ARGV[3]));
		score = n1 + tonumber(ARGV[2]);
	else
		score = tonumber(ARGV[2]);
	end
	if score < 0 then 
		score = 0
	end
	score = score * (2^ARGV[3]) + ARGV[4]
	redis.call('ZADD', KEYS[1], score, ARGV[1]);
	return score;
`)
	currentScore, err = script.Run(context.TODO(), redis_wrapper.Rdb(), []string{redisKey}, member, incrScore, offsetBits, complementTime).Int64()
	return
}

// Sync 同步排行榜
func (rc *RankConf) Sync(regionId int64, member string, score int64) (err error) {
	// 是否同步增量
	//if rc.IsSyncIncr() && !param.IsForceEdit {
	if rc.IsSyncIncr() {
		if rc.IsTimeScore() { // 带时间
			if rc.IsConcurrencySync() { // 是否会并发
				// redis-lua 原子操作
				_, err = RankTimeScoreZIncrByAtomic(rc.RedisKey(regionId), member, score, constants.RankScoreOffsetBits, GetRankCurrentComplementTime())
				if err != nil {
					return err
				}
			} else { // 先获取分数, 处理后再同步
				oldScore, err := rc.GetScore(regionId, member)
				if err != nil {
					return err
				}
				score = CountRankTimeScore(oldScore + score)
				//if score < 0 { // 最低为0分
				//	score = 0
				//}
				err = redis_wrapper.Rdb().ZAdd(context.TODO(), rc.RedisKey(regionId), redis.Z{Score: float64(score), Member: member}).Err()
				if err != nil {
					return err
				}
			}
		} else {
			err = redis_wrapper.Rdb().ZIncrBy(context.TODO(), rc.RedisKey(regionId), float64(score), member).Err()
			if err != nil {
				return err
			}
		}
	} else {
		isNeedSync := true // 是否要同步

		// 分数是否带时间
		if rc.IsTimeScore() {
			// 如果分数不变, 则无需同步
			//if oldScore == 0 {
			//	oldScore, _ = rc.GetScore(regionId, member)
			//}
			//if oldScore == score {
			//	isNeedSync = false
			//}
			//
			//// 降分不同步
			//if rc.DownScoreIsNotSync() {
			//	if score < oldScore {
			//		isNeedSync = false
			//	}
			//}

			score = CountRankTimeScore(score)
		}

		if isNeedSync {
			err = redis_wrapper.Rdb().ZAdd(context.TODO(), rc.RedisKey(regionId), redis.Z{Score: float64(score), Member: member}).Err()
			if err != nil {
				return err
			}
		}
	}

	return
}

// NearbyList 排名附近的列表
func (rc *RankConf) NearbyList(regionId int64, member string, before, after int64) (err error) {
	// 查询当前排名
	currentRank, err := rc.GetRank(regionId, member)
	if err != nil {
		return
	}

	// 获取排行
	start := int64(math.Max(0, float64(currentRank-before)))
	end := currentRank + after
	zl, err := redis_wrapper.Rdb().ZRevRangeWithScores(context.TODO(), rc.RedisKey(regionId), start, end).Result()
	if err != nil {
		return
	}

	for k, v := range zl {
		mem, score, err := ParseUserScoreRet(v, rc.IsTimeScore())
		if err != nil {
			return err
		}

		_ = k
		_ = mem
		_ = score

		//pbRank := &gourmetship.PbRank{
		//	UserID: mem,
		//	Score:  score,
		//	Rank:   int32(k+1) + int32(start),
		//}
		//
		//if userId == uid {
		//	ack.Self = pbRank
		//} else {
		//	if ack.Self == nil {
		//		ack.BeforeList = append(ack.BeforeList, pbRank)
		//	} else {
		//		ack.AfterList = append(ack.AfterList, pbRank)
		//	}
		//}
	}

	return
}

// ScoreNearbyList 分数附近的列表
func (rc *RankConf) ScoreNearbyList(regionId int64, member string, lowPct, highPct float64, count int64) (err error) {
	// 查询分数
	currentScore, err := redis_wrapper.Rdb().ZScore(context.TODO(), rc.RedisKey(regionId), member).Result()
	if err != nil {
		return
	}

	// 获取排行
	min1 := currentScore * lowPct / 100
	max1 := currentScore * highPct / 100
	zl, err := redis_wrapper.Rdb().ZRevRangeByScoreWithScores(context.TODO(), rc.RedisKey(regionId), &redis.ZRangeBy{
		Min: fmt.Sprintf("%f", min1), Max: fmt.Sprintf("%f", max1), Offset: 0, Count: count,
	}).Result()
	if err != nil {
		return
	}

	for _, v := range zl {
		mem, score, err := ParseUserScoreRet(v, rc.IsTimeScore())
		if err != nil {
			return err
		}

		_ = mem
		_ = score
		//pbRank := &gourmetship.PbRank{
		//	UserID: uid,
		//	Score:  score,
		//	//Rank:   int32(k/2 + 1), // 暂不用返回排名, fix排名计算
		//}
		//
		//ack.List = append(ack.List, pbRank)
	}

	return
}
