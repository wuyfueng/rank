package conf

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	// 命令行参数
	flag.String("redis.host", "", "-redis.host")
	flag.Int("redis.port", 0, "-redis.port")
	flag.String("redis.password", "", "-redis.password")

	viper.SetConfigName("conf")    // name of config file (without extension)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./conf/") // path to look for the config file in
	err := viper.ReadInConfig()    // Find and read the config file
	if err != nil {                // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// 加载命令行参数
	flag.Parse()
	// 遍历并设置命令行参数到viper
	flag.Visit(func(f *flag.Flag) {
		viper.Set(f.Name, f.Value.String())
	})
}
