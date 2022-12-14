package locate

import (
	"context"
	"github.com/wha1ever/bfs/internal/rabbitmq"
	"log"
	"os"
	"strconv"
)

func Locate(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate(ctx context.Context) error {

	q := rabbitmq.New(ctx, os.Getenv("RABBITMQ_SERVER"))
	defer func(q *rabbitmq.RabbitMQ) {
		err := q.Close()
		if err != nil {
			log.Println(err)
		}
	}(q)
	err := q.Bind("dataServers")
	if err != nil {
		log.Println(err)
		return err
	}
	c, err := q.Consume()
	for msg := range c {
		object, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			log.Println(err)
			return err
		}
		if Locate(os.Getenv("STORAGE_ROOT") + "/object/" + object) {
			err := q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
			if err != nil {
				log.Println(err)
				return err
			}
		}
	}
	return nil
}
