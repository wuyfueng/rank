package rank

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/wuyfueng/rank/common/proto"
	"github.com/wuyfueng/tools"
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
	if rc.IsDense() { // 密集排名
		// 先获取分数
		score, err := rc.GetScore(regionId, member)
		if err != nil {
			return 0, err
		}

		// 从 RedisDenseKey 获取
		if rc.IsPositiveSort() { // 是否正序
			rank, err = redis_wrapper.Rdb().ZRank(context.TODO(), rc.RedisDenseKey(regionId), fmt.Sprintf("%d", score)).Result()
			if err != nil {
				return 0, err
			}
		} else {
			rank, err = redis_wrapper.Rdb().ZRevRank(context.TODO(), rc.RedisDenseKey(regionId), fmt.Sprintf("%d", score)).Result()
			if err != nil {
				return 0, err
			}
		}
	} else {                     // 正常排名
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
	var oldScore, newScore int64 // 不带时间戳的新旧分数
	// 密集排名, 先获取旧分数
	if rc.IsDense() {
		oldScore, err = rc.GetScore(regionId, member)
		if err != nil {
			return err
		}
	}

	if rc.IsSyncIncr() { // 增量同步
		if rc.IsTimeScore() { // 带时间
			if rc.IsConcurrencySync() { // 是否会并发
				// redis-lua 原子操作
				currentScore, err := RankTimeScoreZIncrByAtomic(rc.RedisKey(regionId), member, score, constants.RankScoreOffsetBits, GetRankCurrentComplementTime())
				if err != nil {
					return err
				}
				newScore = GetRankRealScore(currentScore)
			} else { // 先获取分数, 处理后再同步
				if oldScore == 0 {
					oldScore, err = rc.GetScore(regionId, member)
					if err != nil {
						return err
					}
				}
				newScore = oldScore + score
				score = CountRankTimeScore(newScore)
				//if score < 0 { // 最低为0分
				//	score = 0
				//}
				err = redis_wrapper.Rdb().ZAdd(context.TODO(), rc.RedisKey(regionId), redis.Z{Score: float64(score), Member: member}).Err()
				if err != nil {
					return err
				}
			}
		} else {
			fNewScore, err := redis_wrapper.Rdb().ZIncrBy(context.TODO(), rc.RedisKey(regionId), float64(score), member).Result()
			if err != nil {
				return err
			}
			newScore = int64(fNewScore)
		}
	} else { // 覆盖同步
		newScore = score

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

		if !isNeedSync {
			return
		}
		err = redis_wrapper.Rdb().ZAdd(context.TODO(), rc.RedisKey(regionId), redis.Z{Score: float64(score), Member: member}).Err()
		if err != nil {
			return err
		}
	}

	// 密集排名
	if rc.IsDense() {
		// 旧分数是否还存在
		isExistOld := false
		if rc.IsTimeScore() { // 带时间戳
			zl, err := redis_wrapper.Rdb().ZRangeByScore(context.TODO(), rc.RedisKey(regionId), &redis.ZRangeBy{
				Min: fmt.Sprintf("%d", score<<constants.RankScoreOffsetBits), Max: fmt.Sprintf("%d", (score+1)<<constants.RankScoreOffsetBits-1), Offset: 0, Count: 1,
			}).Result()
			if err == nil && len(zl) > 0 {
				isExistOld = true
			}
		} else { // 不带时间戳
			zl, err := redis_wrapper.Rdb().ZRangeByScore(context.TODO(), rc.RedisKey(regionId), &redis.ZRangeBy{
				Min: fmt.Sprintf("%d", oldScore), Max: fmt.Sprintf("%d", oldScore), Offset: 0, Count: 1,
			}).Result()
			if err == nil && len(zl) > 0 {
				isExistOld = true
			}
		}
		// 旧分数不存在, 从 RedisDenseKey 移除
		if !isExistOld {
			err = redis_wrapper.Rdb().ZRem(context.TODO(), rc.RedisDenseKey(regionId), oldScore).Err()
			if err != nil {
				// err_log
			}
		}

		// 添加新分数到 RedisDenseKey
		err = redis_wrapper.Rdb().ZAdd(context.TODO(), rc.RedisDenseKey(regionId), redis.Z{Score: float64(newScore), Member: newScore}).Err()
		if err != nil {
			// err_log
		}
	}

	return
}

// TopList 排行榜top列表
func (rc *RankConf) TopList(regionId int64, n int64) (list []*pb.PbRank, err error) {
	// 上榜人数
	if n <= 0 {
		n = int64(rc.RankingUserNum())
	}
	if n <= 0 {
		return nil, errors.New(fmt.Sprintf("TopList n: %d is invalid", n))
	}

	if rc.IsDense() { // 密集排名
		// 获取第n名的分数
		scoreZl, err := redis_wrapper.Rdb().ZRevRangeWithScores(context.TODO(), rc.RedisDenseKey(regionId), n-1, n-1).Result()
		if err != nil {
			return nil, err
		}
		if len(scoreZl) != 1 {
			return nil, errors.New("scoreZl is empty")
		}

		zl := make([]redis.Z, 0, n)

		// 根据分数范围查找
		minScore := scoreZl[0].Score // 不带时间戳的分数下限
		if rc.IsTimeScore() {        // 带时间戳
			zl, err = redis_wrapper.Rdb().ZRevRangeByScoreWithScores(context.TODO(), rc.RedisKey(regionId), &redis.ZRangeBy{
				Min: fmt.Sprintf("%d", int64(minScore)<<constants.RankScoreOffsetBits), Max: "+inf", Offset: 0, Count: -1,
			}).Result()
			if err != nil {
				return nil, err
			}
		} else { // 不带时间戳
			zl, err = redis_wrapper.Rdb().ZRevRangeByScoreWithScores(context.TODO(), rc.RedisKey(regionId), &redis.ZRangeBy{
				Min: fmt.Sprintf("%f", minScore), Max: "+inf", Offset: 0, Count: -1,
			}).Result()
			if err != nil {
				return nil, err
			}
		}

		lastScore := int64(math.MinInt)
		curRank := int32(0)
		for _, v := range zl {
			member, score, err := ParseUserScoreRet(v, rc.IsTimeScore())
			if err != nil {
				return nil, err
			}

			if score != lastScore {
				lastScore = score
				curRank++
			}
			list = append(list, &pb.PbRank{
				Member: member,
				Score:  score,
				Rank:   curRank,
			})
		}
	} else { // 正常排名
		zl := make([]redis.Z, 0, n)
		if rc.IsPositiveSort() { // 是否正序
			zl, err = redis_wrapper.Rdb().ZRangeWithScores(context.TODO(), rc.RedisKey(regionId), 0, n-1).Result()
			if err != nil {
				return nil, err
			}
		} else {
			zl, err = redis_wrapper.Rdb().ZRevRangeWithScores(context.TODO(), rc.RedisKey(regionId), 0, n-1).Result()
			if err != nil {
				return nil, err
			}
		}

		for k, v := range zl {
			member, score, err := ParseUserScoreRet(v, rc.IsTimeScore())
			if err != nil {
				return nil, err
			}

			list = append(list, &pb.PbRank{
				Member: member,
				Score:  score,
				Rank:   int32(k + 1),
			})
		}
	}

	return
}

// NearbyList 排名附近的列表
func (rc *RankConf) NearbyList(regionId int64, member string, before, after int64) (list []*pb.PbRank, err error) {
	// 查询当前排名
	currentRank, err := rc.GetRank(regionId, member)
	if err != nil {
		return nil, err
	}

	// 获取排行
	currentRank--
	start := int64(math.Max(0, float64(currentRank-before)))
	end := currentRank + after
	if rc.IsDense() { // 密集排名
		// 先根据名次获取分数范围
		scoreZl, err := redis_wrapper.Rdb().ZRevRangeWithScores(context.TODO(), rc.RedisDenseKey(regionId), start, end).Result()
		if err != nil {
			return nil, err
		}
		if len(scoreZl) == 0 {
			return nil, errors.New("NearbyList RedisDenseKey not found")
		}

		zl := make([]redis.Z, 0, 1)

		// 根据分数范围查找
		if rc.IsTimeScore() { // 带时间戳
			zl, err = redis_wrapper.Rdb().ZRevRangeByScoreWithScores(context.TODO(), rc.RedisKey(regionId), &redis.ZRangeBy{
				Min: fmt.Sprintf("%d", int64(scoreZl[len(scoreZl)-1].Score)<<constants.RankScoreOffsetBits), Max: fmt.Sprintf("%d", int64(scoreZl[0].Score+1)<<constants.RankScoreOffsetBits-1), Offset: 0, Count: -1,
			}).Result()
			if err != nil {
				return nil, err
			}
		} else { // 不带时间戳
			zl, err = redis_wrapper.Rdb().ZRevRangeByScoreWithScores(context.TODO(), rc.RedisKey(regionId), &redis.ZRangeBy{
				Min: fmt.Sprintf("%f", scoreZl[len(scoreZl)-1].Score), Max: fmt.Sprintf("%f", scoreZl[0].Score), Offset: 0, Count: -1,
			}).Result()
			if err != nil {
				return nil, err
			}
		}

		lastScore := int64(math.MinInt)
		curRank := int32(start)
		for _, v := range zl {
			mem, score, err := ParseUserScoreRet(v, rc.IsTimeScore())
			if err != nil {
				return nil, err
			}

			if lastScore != score {
				lastScore = score
				curRank++
			}
			pbRank := &pb.PbRank{
				Member: mem,
				Score:  score,
				Rank:   curRank,
			}

			list = append(list, pbRank)
		}
	} else { // 正常情况
		zl, err := redis_wrapper.Rdb().ZRevRangeWithScores(context.TODO(), rc.RedisKey(regionId), start, end).Result()
		if err != nil {
			return nil, err
		}
		for k, v := range zl {
			mem, score, err := ParseUserScoreRet(v, rc.IsTimeScore())
			if err != nil {
				return nil, err
			}

			pbRank := &pb.PbRank{
				Member: mem,
				Score:  score,
				Rank:   int32(k+1) + int32(start),
			}

			list = append(list, pbRank)
		}
	}

	return
}

// ScoreNearbyList 分数附近的列表
func (rc *RankConf) ScoreNearbyList(regionId int64, member string, lowPct, highPct float64, count int64) (list []*pb.PbRank, err error) {
	// 查询分数
	currentScore, err := redis_wrapper.Rdb().ZScore(context.TODO(), rc.RedisKey(regionId), member).Result()
	if err != nil {
		return nil, err
	}

	// 获取排行
	min1 := currentScore * lowPct / 100
	max1 := currentScore * highPct / 100
	zl, err := redis_wrapper.Rdb().ZRevRangeByScoreWithScores(context.TODO(), rc.RedisKey(regionId), &redis.ZRangeBy{
		Min: fmt.Sprintf("%f", min1), Max: fmt.Sprintf("%f", max1), Offset: 0, Count: count,
	}).Result()
	if err != nil {
		return nil, err
	}

	for _, v := range zl {
		mem, score, err := ParseUserScoreRet(v, rc.IsTimeScore())
		if err != nil {
			return nil, err
		}

		pbRank := &pb.PbRank{
			Member: mem,
			Score:  score,
		}

		list = append(list, pbRank)
	}

	return
}

// CreateDenseData 创建密集排名
func (rc *RankConf) CreateDenseData(regionId int64) (err error) {
	count, err := redis_wrapper.Rdb().ZCard(context.TODO(), rc.RedisKey(regionId)).Result()
	if err != nil {
		return err
	}

	// 分页处理
	tools.PageRange(int(count), 1000, func(page, startIndex, endIndex int) {
		zl, err := redis_wrapper.Rdb().ZRangeWithScores(context.TODO(), rc.RedisKey(regionId), int64(startIndex), int64(endIndex)).Result()
		if err != nil {
			// err_log
			return
		}

		m := make(map[int64]struct{}, 1)
		for _, v := range zl {
			score := int64(v.Score)
			if rc.IsTimeScore() { // 去除时间戳
				score = GetRankRealScore(score)
			}
			if _, ok := m[score]; !ok {
				m[score] = struct{}{}
			}
		}
		for k := range m {
			err = redis_wrapper.Rdb().ZAdd(context.TODO(), rc.RedisDenseKey(regionId), redis.Z{Score: float64(k), Member: k}).Err()
			if err != nil {
				// err_log
			}
		}
	})

	return
}
