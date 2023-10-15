package tool

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
	"sync"
)

var HaloTail *tail.Tail
var SqlTail *tail.Tail

// InitTail 初始化Tail
func InitTail() {
	HaloFile := viper.GetString("HaloLog.path")
	SqlFile := viper.GetString("SqlLog.path")
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	// 打开文件开始读取数据
	halotails, err := tail.TailFile(HaloFile, config)
	if err != nil {
		logrus.Error("tail %s failed, err:%v\n\n", HaloFile, err)
		return
	}

	sqltails, err := tail.TailFile(SqlFile, config)
	if err != nil {
		logrus.Error("tail %s failed, err:%v\n\n", HaloFile, err)
		return
	}

	HaloTail = halotails
	SqlTail = sqltails
	return
}

// ReadLog 读取每一行日志丢到通道
func ReadLog(wg *sync.WaitGroup) {
	wg.Add(2)
	go ReadHaloLog(wg)
	go ReadSqlLog(wg)

	//循环读数据
	//for {
	//	line, ok := <-HaloTail.Lines
	//	if !ok {
	//		logrus.Warn("tail file close , filename:", HaloTail.Filename)
	//		continue
	//	}
	//	// 如果是空行就跳过
	//	if len(strings.Trim(line.Text, "\r")) == 0 {
	//		continue
	//	}
	//	// 测试是否能拿到msg
	//	fmt.Println("msg:", line.Text)
	//
	//	// 把读出来的每一行数据包装成msg类型,发送到Kafka
	//	msg := &sarama.ProducerMessage{
	//		Topic: viper.GetString("HaloLog.topic"),
	//		Value: sarama.StringEncoder(line.Text),
	//	}
	//	// 包装完成之后丢到通道中
	//	HaloLogMsgChan <- msg
	//	fmt.Println("消息发送成功")
	//}
}

func ReadHaloLog(wg *sync.WaitGroup) {
	//循环读数据
	for {
		line, ok := <-HaloTail.Lines
		if !ok {
			logrus.Warn("tail file close , filename:", HaloTail.Filename)
			continue
		}
		// 如果是空行就跳过
		if len(strings.Trim(line.Text, "\r")) == 0 {
			continue
		}
		// 测试是否能拿到msg
		fmt.Println("halo.log msg:", line.Text)

		// 把读出来的每一行数据包装成msg类型,发送到Kafka
		msg := &sarama.ProducerMessage{
			Topic: viper.GetString("HaloLog.topic"),
			Value: sarama.StringEncoder(line.Text),
		}
		// 包装完成之后丢到通道中
		HaloLogMsgChan <- msg
		fmt.Println("消息发送成功")
	}
}

func ReadSqlLog(wg *sync.WaitGroup) {
	//循环读数据
	for {
		line, ok := <-SqlTail.Lines
		if !ok {
			logrus.Warn("tail file close , filename:", SqlTail.Filename)
			continue
		}
		// 如果是空行就跳过
		if len(strings.Trim(line.Text, "\r")) == 0 {
			continue
		}
		// 测试是否能拿到msg
		fmt.Println("sql.log msg:", line.Text)

		// 把读出来的每一行数据包装成msg类型,发送到Kafka
		//msg := &sarama.ProducerMessage{
		//	Topic: viper.GetString("SqlLog.topic"),
		//	Value: sarama.StringEncoder(line.Text),
		//}
		// 包装完成之后丢到通道中
		//SqlLogMsgChan <- msg
	}
}
