package main

import (
	"LogAgent/tool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

func main() {
	// 加载配置
	tool.LoadConfig()
	logrus.Info("Init Config")

	//初始化etcd
	tool.InitEtcd()

	//从etcd拉取最新的配置
	tool.GetConfig(viper.GetString("etcd.collect_key"))

	// 初始化tail
	tool.InitTail()

	//初始化Kafka
	tool.InitKafka()

	wg := &sync.WaitGroup{}
	// 读取每一行日志
	tool.ReadLog(wg)

	wg.Wait()
}
