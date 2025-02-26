package demo

import (
	"fmt"
	"log"
	"testing"
)

var playerId = "1"

// 创建测试数据
func TestCreateTestData(t *testing.T) {
	createTestData()
}

// 更新玩家积分
func TestUpdateScore(t *testing.T) {
	updateScore(playerId, 100)
}

// 获取玩家当前排名
func TestGetPlayerRank(t *testing.T) {
	rank, _ := getPlayerRank(playerId)
	log.Printf("玩家: %s, 当前排名: %d", playerId, rank)
}

// 获取排行榜前N名
func TestGetTopN(t *testing.T) {
	n := int64(10)
	list, _ := getTopN(n)
	log.Printf("前: %d名列表如下", n)
	for _, v := range list {
		fmt.Println(fmt.Sprintf("第: %d名, playerId: %s, score: %d", v.Rank, v.Member, v.Score))
	}
}

// 获取玩家周边排名
func TestGetPlayerRankRange(t *testing.T) {
	before := int64(2)
	after := int64(3)
	log.Printf("玩家id: %s, 周边前: %d名, 后: %d名列表如下", playerId, before, after)
	list, _ := getPlayerRankRange(playerId, before, after)
	for _, v := range list {
		fmt.Println(fmt.Sprintf("第: %d名, playerId: %s, score: %d", v.Rank, v.Member, v.Score))
	}
}

// 创建密集排名
func TestCreateDenseData(t *testing.T) {
	createDenseData()
}
