package tool

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var Client *clientv3.Client

func InitEtcd() {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{viper.GetString("etcd.addr")},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		log.Fatal("etcd link fail: ", err)
	}
	Client = client

}

// 拉取日志配置收集配置项的函数
func GetConfig(key string) {
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := Client.Get(ctx, key)
	defer cancel()
	if err != nil {
		logrus.Errorf("get from etcd failed, err:%v\n", err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warnf("can't get any value by key:%s from etcd", key)
		return
	}
	//keyValues := resp.Kvs[0]
	// json格式字符串
	//err = json.Unmarshal(keyValues.Value, &conf)
	//if err != nil {
	//	logrus.Errorf("unmarshal value from etcd failed, err:%v", err)
	//	return
	//}
	//log.Debugf("load conf from etcd success, conf:%#v", conf)
	//return
}
