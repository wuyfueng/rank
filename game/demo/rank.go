package demo

import (
	"context"
	"fmt"
	"time"

	"github.com/wuyfueng/rank/common/constants"
	"github.com/wuyfueng/rank/common/rank"
	"github.com/wuyfueng/rank/common/redis_wrapper"
)

// Sync 更新玩家积分
func Sync() {
	c := rank.GetRankConf(constants.RankTypeScore)
	if c == nil {
		fmt.Println("GetRankConf is nil")
		return
	}

	fmt.Println(c.Name())

	str, err := redis_wrapper.Rdb().Set(context.TODO(), "yf", 33, time.Second*109).Result()
	fmt.Println("err", err)
	fmt.Println("str", str)
	str1, err := redis_wrapper.Rdb().Get(context.TODO(), "yf").Result()
	fmt.Println("err1", err)
	fmt.Println("str1", str1)
}
