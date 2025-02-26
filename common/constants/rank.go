package constants

import "time"

type RankRegionType int // 排行榜分区(范围)类型
const (
	RankRegionTypeGlobal RankRegionType = 1 // 全服
	RankRegionTypeServer RankRegionType = 2 // 区服
)

type RankType int // 排行榜类型
const (
	RankTypeScore RankType = 1 // 全服分数排行榜
)

const (
	// RankFinalTime - 当前时间 用于追加分数, 截止时间 2036-1-1 00:00:00
	RankFinalTime = 2082729600
	// RankScoreOffsetBits 分数左移位数 分数上限 8589934592
	RankScoreOffsetBits = 29

	RankCacheDurationSize = time.Second * 60 // 缓存时长最小粒度
)
