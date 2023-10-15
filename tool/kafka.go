package tool

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Producer       sarama.SyncProducer
	HaloLogMsgChan chan *sarama.ProducerMessage
	SqlLogMsgChan  chan *sarama.ProducerMessage
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

	HaloLogMsgChan = make(chan *sarama.ProducerMessage, viper.GetInt64("Kafka.msg_chan"))

	SqlLogMsgChan = make(chan *sarama.ProducerMessage, viper.GetInt64("Kafka.msg_chan"))

	Producer = producer

	go SendMsg()
}

// SendMsg 发送日志中的数据到Kafka
func SendMsg() {
	for {
		select {
		case msg := <-HaloLogMsgChan:
			partition, offset, err := Producer.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msg failed , err:", err)
			}
			logrus.Info("send msg to Kafka success: pid: ", partition, " offset: ", offset)
		}
	}
}
