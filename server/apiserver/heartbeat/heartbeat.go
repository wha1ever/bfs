package heartbeat

import (
	"context"
	"github.com/wha1ever/bfs/internal/rabbitmq"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	//dataServers = make(map[string]time.Time)
	//mutex       sync.Mutex
	dataServers sync.Map
)

func ListenHeartbeat(ctx context.Context) error {
	q := rabbitmq.New(ctx, os.Getenv("RABBITMQ_SERVER"))
	defer func(q *rabbitmq.RabbitMQ) {
		err := q.Close()
		if err != nil {
			log.Println(err)
		}
	}(q)
	err := q.Bind("apiServers")
	if err != nil {
		log.Println(err)
		return err
	}
	c, err := q.Consume()
	if err != nil {
		log.Println(err)
		return err
	}
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			log.Println(err)
			return err
		}
		//mutex.Lock()
		//dataServers[dataServer] = time.Now()
		//mutex.Unlock()
		dataServers.Store(dataServer, time.Now())
	}
	return nil
}

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		dataServers.Range(func(key, value any) bool {
			if value.(time.Time).Add(10 * time.Second).Before(time.Now()) { // 10s没有收到心跳包就会从map中删除
				dataServers.Delete(key)
			}
			return true
		})
	}
}

func GetDataServers() []string {
	ds := make([]string, 0)
	dataServers.Range(func(key, value any) bool {
		ds = append(ds, key.(string))
		return true
	})
	return ds
}

// ChooseRandomDataServer Todo: 负载均衡实现
func ChooseRandomDataServer() string {
	ds := GetDataServers()
	n := len(ds)
	if n == 0 {
		return ""
	}
	return ds[rand.Intn(n)]
}
