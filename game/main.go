package main

import (
	_ "github.com/wuyfueng/rank/game/conf"

	"github.com/wuyfueng/rank/common/redis_wrapper"
	"github.com/wuyfueng/tools"

	"github.com/wuyfueng/rank/game/demo"
)

func main() {
	// 初始化redis
	redis_wrapper.Init()

	demo.Sync()

	tools.WaitExit()
}
