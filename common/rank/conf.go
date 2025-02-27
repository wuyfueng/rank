package rank

import (
	"fmt"
	"time"

	"github.com/wuyfueng/rank/common/constants"
)

// RankConf 排行榜配置
type RankConf struct {
	rankType          constants.RankType
	desc              string                   // 描述
	name              string                   // 名称
	regionType        constants.RankRegionType // 分区类型
	redisKey          string                   // redisKey
	likeRedisKey      string                   // 点赞redisKey(hash k-v:userId-likes)
	notifyRedisKey    string                   // 登陆通知CD
	isHomePage        bool                     // 是否在排行榜首页
	isPositiveSort    bool                     // 是否正序
	isTimeScore       bool                     // 分数是否带时间
	isSyncIncr        bool                     // 是否同步增量
	isConcurrencySync bool                     // 是否会并发同步
	cacheUserNum      int                      // 缓存人数
	rankingUserNum    int                      // 上榜人数
	cacheDuration     time.Duration            // 缓存时长
	isTimely          bool                     // 是否是实时获取排行
	isDense           bool                     // 是否是密集排名
}

// 注册排行榜配置
var (
	// 类型对应配置
	rankConfMap = map[constants.RankType]*RankConf{
		constants.RankTypeScore: {rankType: constants.RankTypeScore, desc: "全服分数排行榜", name: "score", regionType: constants.RankRegionTypeGlobal, isHomePage: true, isPositiveSort: false, isTimeScore: true, isSyncIncr: false, isConcurrencySync: false, cacheUserNum: 100, rankingUserNum: 100, cacheDuration: constants.RankCacheDurationSize, isTimely: true, isDense: false},
	}
)

func RankConfList() map[constants.RankType]*RankConf {
	return rankConfMap
}

func GetRankConf(rankType constants.RankType) *RankConf {
	return rankConfMap[rankType]
}

func (rc *RankConf) RankType() constants.RankType         { return rc.rankType }
func (rc *RankConf) Desc() string                         { return rc.desc }
func (rc *RankConf) Name() string                         { return rc.name }
func (rc *RankConf) RegionType() constants.RankRegionType { return rc.regionType }
func (rc *RankConf) IsHomePage() bool                     { return rc.isHomePage }
func (rc *RankConf) IsPositiveSort() bool                 { return rc.isPositiveSort }
func (rc *RankConf) IsTimeScore() bool                    { return rc.isTimeScore }
func (rc *RankConf) IsSyncIncr() bool                     { return rc.isSyncIncr }
func (rc *RankConf) IsConcurrencySync() bool              { return rc.isConcurrencySync }
func (rc *RankConf) CacheUserNum() int                    { return rc.cacheUserNum }
func (rc *RankConf) RankingUserNum() int                  { return rc.rankingUserNum }
func (rc *RankConf) CacheDuration() time.Duration         { return rc.cacheDuration }
func (rc *RankConf) IsNeedCache() bool                    { return rc.cacheDuration > 0 } // 是否需要缓存
func (rc *RankConf) IsTimely() bool                       { return rc.isTimely }
func (rc *RankConf) IsDense() bool                        { return rc.isDense }

func (rc *RankConf) spellRedisKey(regionId int64, keyName string) string {
	switch rc.regionType {
	case constants.RankRegionTypeGlobal:
		return fmt.Sprintf("global:%s:%s", keyName, rc.name)
	case constants.RankRegionTypeServer:
		return fmt.Sprintf("server_%04d:%s:%s", regionId, keyName, rc.name)
	}

	// err_log
	return ""
}
func (rc *RankConf) RedisKey(regionId int64) (key string) {
	return rc.spellRedisKey(regionId, "Rank")
}
func (rc *RankConf) RedisDenseKey(regionId int64) (key string) {
	return rc.spellRedisKey(regionId, "RankDense")
}
func (rc *RankConf) LikeRedisKey(regionId int64) (key string) {
	return rc.spellRedisKey(regionId, "RankLike")
}
func (rc *RankConf) NotifyRedisKey(regionId int64, notifykey string) (key string) {
	return fmt.Sprintf("%s:%s", rc.spellRedisKey(regionId, "RankNotify"), notifykey)
}
func (rc *RankConf) WaitRedisKey() string {
	return fmt.Sprintf("RankWait:%s", rc.name)
}

// CountRankTimeScore 计算排行榜追加时间后的分数
func CountRankTimeScore(score int64) int64 {
	return score<<constants.RankScoreOffsetBits + GetRankCurrentComplementTime()
}

// GetRankRealScore 获取排行榜真实分数
func GetRankRealScore(score int64) int64 {
	return score >> constants.RankScoreOffsetBits
}

// GetRankCurrentComplementTime 获取排行榜当前时间补时
func GetRankCurrentComplementTime() int64 {
	return constants.RankFinalTime - time.Now().Unix()
}
