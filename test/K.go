package test

import (
	"LogAgent/tool"
)

func Test() {
	tool.LoadConfig()

	tool.InitKafka()

}
