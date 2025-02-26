package main

import (
	"github.com/wuyfueng/rank/common/redis_wrapper"
	_ "github.com/wuyfueng/rank/game/conf"
	"github.com/wuyfueng/tools"
)

func main() {
	// 初始化redis
	redis_wrapper.Init()

	tools.WaitExit()
}
