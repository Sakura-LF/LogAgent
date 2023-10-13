package tool

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var Tail *tail.Tail

// InitTail 初始化Tail
func InitTail() {
	file := viper.GetString("log_file.halo_log.path")
	config := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}

	// 打开文件开始读取数据
	tails, err := tail.TailFile(file, config)
	if err != nil {
		logrus.Error("tail %s failed, err:%v\n\n", file, err)
		return
	}
	Tail = tails
	return
}

// ReadLog 读取每一行日志丢到通道
func ReadLog() {
	//循环读数据
	for {
		line, ok := <-Tail.Lines
		if !ok {
			logrus.Warn("tail file close , filename:", Tail.Filename)
			continue
		}
		// 如果是空行就跳过
		if len(strings.Trim(line.Text, "\r")) == 0 {
			continue
		}
		// 测试是否能拿到msg
		fmt.Println("msg:", line.Text)

		// 把读出来的每一行数据包装成msg类型,发送到Kafka
		msg := &sarama.ProducerMessage{
			Topic: viper.GetString("log_file.halo_log.topic"),
			Value: sarama.StringEncoder(line.Text),
		}
		// 包装完成之后丢到通道中
		MsgChan <- msg
	}
}
