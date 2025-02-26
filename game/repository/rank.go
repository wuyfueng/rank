package demo

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/wuyfueng/rank/common/constants"
	"github.com/wuyfueng/rank/common/rank"
)

// Sync 更新玩家积分
func Sync() {
	rc := rank.GetRankConf(constants.RankTypeScore)
	if rc == nil {
		fmt.Println("GetRankConf is nil")
		return
	}

	for i := 1; i <= 200; i++ {
		err := rc.Sync(0, fmt.Sprintf("%d", i), rand.Int63n(100))
		if err != nil {
			log.Println("Sync err", err)
		}
	}
}
