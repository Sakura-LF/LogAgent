package main

import (
	"LogAgent/tool"
	"github.com/sirupsen/logrus"
)

func main() {
	// 加载配置
	tool.LoadConfig()
	logrus.Info("Init Config")

	//初始化etcd
	tool.InitEtcd()

	//从etcd拉取最新的配置
	//tool.GetConfig()

	// 初始化tail
	tool.InitTail()

	//初始化Kafka
	tool.InitKafka()

	// 读取每一行日志
	tool.ReadLog()
}
