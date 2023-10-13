package tool

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Producer sarama.SyncProducer
	MsgChan  chan *sarama.ProducerMessage
)

// InitKafka 初始化Kafka
func InitKafka() {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{viper.GetString("Kafka.addr")}, config)
	if err != nil {
		logrus.Errorf("producer closed, err:", err)
		return
	}
	MsgChan = make(chan *sarama.ProducerMessage, viper.GetInt64("Kafka.msg_chan"))
	Producer = producer

	go SendMsg()

	//1.生产者配置
	//config := sarama.NewConfig()
	//config.Producer.RequiredAcks = sarama.WaitForAll          //发送完数据需要leader和follow都确认
	//config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	//config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	//
	//// 2.连接Kafka
	//producer, err := sarama.NewSyncProducer([]string{viper.GetString("Kafka.addr")}, nil)
	//if err != nil {
	//	logrus.Error("link fail", err)
	//}
	//
	//// 定义存放消息的管道,异步通信
	//MsgChan = make(chan *sarama.ProducerMessage, viper.GetInt64("Kafka.msg_chan"))
	//
	//Producer = producer
	//
	//go SendMsg()
}

// SendMsg 发送日志中的数据到Kafka
func SendMsg() {
	for {
		select {
		case msg := <-MsgChan:
			partition, offset, err := Producer.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msg failed , err:", err)
			}
			logrus.Info("send msg to Kafka success: pid: ", partition, " offset: ", offset)
		}
	}
}
