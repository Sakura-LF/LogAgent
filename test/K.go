package test

import (
	"LogAgent/tool"
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

func Test() {
	tool.LoadConfig()

	tool.InitKafka()

}

func Etcd() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		log.Fatal("etcd link fail: ", err)
	}

	defer cli.Close()

	//// put
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//_, err = cli.Put(ctx, "name", "Sakura")
	//cancel()
	//if err != nil {
	//	fmt.Println("Put failed:", err)
	//}
	//
	//// get
	//ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	//gr, err := cli.Get(ctx, "name")
	//cancel()
	//if err != nil {
	//	fmt.Println("get failed :", err)
	//	return
	//}

	////
	//for _, value := range gr.Kvs {
	//	fmt.Println("key: ", string(value.Key))
	//	fmt.Println("value: ", string(value.Value))
	//}

	watch := cli.Watch(context.Background(), "name")
	for v := range watch {
		for _, value := range v.Events {
			fmt.Println("type:", value.Type, " key: ", string(value.Kv.Key), " value: ", string(value.Kv.Value))
		}
	}

}

func Etcd2() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		log.Fatal("etcd link fail: ", err)
	}

	defer cli.Close()

	// put
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = cli.Put(ctx, "name", "Sakura1")
	cancel()
	if err != nil {
		fmt.Println("Put failed:", err)
	}

	//// get
	//ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	//gr, err := cli.Get(ctx, "name")
	//cancel()
	//if err != nil {
	//	fmt.Println("get failed :", err)
	//	return
	//}
	//
	//////
	////for _, value := range gr.Kvs {
	////	fmt.Println("key: ", string(value.Key))
	////	fmt.Println("value: ", string(value.Value))
	////}
}
