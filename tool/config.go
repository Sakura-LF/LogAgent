package tool

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Task struct {
	HaloLog `mapstructure:"HaloLog"`
	SqlLog  `mapstructure:"SqlLog"`
}

type HaloLog struct {
	Topic string `mapstructure:"topic"`
	Path  string `mapstructure:"path"`
}
type SqlLog struct {
	Topic string `mapstructure:"topic"`
	Path  string `mapstructure:"path"`
}

func LoadConfig() {
	// 定义默认值,配置文件加载错误使用的值
	SetDefault()

	// 配置文件属性设置
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	//读取配置
	if err := viper.ReadInConfig(); err != nil {
		logrus.Error(err)
	}
}

func SetDefault() {
	viper.SetDefault("Kafka.addr", 9092)
	viper.SetDefault("log_file.addr", "/halo/halo2/logs/halo.log")
	viper.SetDefault("Kafka.msg_chan", 10000)
}
