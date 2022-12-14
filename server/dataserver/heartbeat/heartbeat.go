package heartbeat

import (
	"context"
	"github.com/wha1ever/bfs/internal/rabbitmq"
	"log"
	"os"
	"time"
)

//
// StartHeartbeat
//  @Description: 每5s发送心跳包给apiServers exchange，内容是本节点的监听地址
//
func StartHeartbeat() {
	ctx := context.Background()
	q := rabbitmq.New(ctx, os.Getenv("RABBITMQ_SERVER"))
	defer func(q *rabbitmq.RabbitMQ) {
		err := q.Close()
		if err != nil {
			log.Println(err)
		}
	}(q)
	for {
		err := q.Publish("apiServer", os.Getenv("LISTEN_ADDRESS"))
		if err != nil {
			log.Println(err)
		}
		time.Sleep(5 * time.Second)
	}
}
