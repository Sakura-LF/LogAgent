package test

import (
	"LogAgent/tool"
	"context"
	"fmt"
	"github.com/IBM/sarama"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"sync"
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

func KafkaConsumer() {
	// 创建新的消费者
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	// 拿到指定topic下面的所有分区列表
	partitionList, err := consumer.Partitions("HaloLogs") // 根据topic取到所有的分区
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Println(partitionList)
	var wg sync.WaitGroup
	for partition := range partitionList { // 遍历所有的分区
		// 针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("HaloLogs", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n",
				partition, err)
			return
		}
		defer pc.AsyncClose()
		// 异步从每个分区消费信息
		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%s Value:%s",
					msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}
	wg.Wait()
}
