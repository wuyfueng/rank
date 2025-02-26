package demo

import (
	"log"
	
	"github.com/wuyfueng/rank/common/constants"
	pb "github.com/wuyfueng/rank/common/proto"
	"github.com/wuyfueng/rank/common/rank"
	"github.com/wuyfueng/rank/common/redis_wrapper"
)

var (
	rc *rank.RankConf // 要测试的排行榜
)

func init() {
	redis_wrapper.RegisterRdb("127.0.0.1", 6379, "")

	rc = rank.GetRankConf(constants.RankTypeScore)
	if rc == nil {
		log.Panic("GetRankConf is nil")
		return
	}
}

// 更新玩家积分
func updateScore(playerId string, score int64) {
	err := rc.Sync(0, playerId, score)
	if err != nil {
		log.Printf("updateScore playerId: %s, score: %d, err: %v", playerId, score, err)
	}
}

// 获取玩家当前排名
func getPlayerRank(playerId string) (rank int64, err error) {
	rank, err = rc.GetRank(0, playerId)
	if err != nil {
		log.Printf("getPlayerRank playerId: %s, err: %v", playerId, err)
	}
	return
}

// 获取排行榜前N名
func getTopN(n int64) (list []*pb.PbRank, err error) {
	list, err = rc.TopList(0, n)
	if err != nil {
		log.Printf("getTopN n: %d, err: %v", n, err)
	}
	return
}

// 获取玩家周边排名
func getPlayerRankRange(playerId string, before, after int64) (list []*pb.PbRank, err error) {
	list, err = rc.NearbyList(0, playerId, before, after)
	if err != nil {
		log.Printf("getPlayerRankRange playerId: %s, before: %d, after: %d, err: %v", playerId, before, after, err)
	}
	return
}
