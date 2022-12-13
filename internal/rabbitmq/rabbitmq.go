package rabbitmq

import (
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

type RabbitMQ struct {
	ctx      context.Context
	channel  *amqp.Channel
	Name     string
	exchange string
}

func New(ctx context.Context, url string) *RabbitMQ {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Println(err)
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		return nil
	}
	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
		return nil
	}
	mq := new(RabbitMQ)
	mq.ctx = ctx
	mq.channel = ch
	mq.Name = q.Name
	return mq
}

func (q *RabbitMQ) Bind(exchange string) error {
	err := q.channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	q.exchange = exchange
	return nil
}

func (q *RabbitMQ) Send(queue string, body interface{}) error {
	str, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return err
	}
	err = q.channel.PublishWithContext(q.ctx, "", queue, false, false, amqp.Publishing{
		ReplyTo:   q.Name,
		Timestamp: time.Now(),
		Body:      []byte(str),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (q *RabbitMQ) Publish(exchange string, body interface{}) error {
	str, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		return err
	}
	err = q.channel.PublishWithContext(q.ctx, q.exchange, "", false, false, amqp.Publishing{
		ReplyTo:   q.Name,
		Timestamp: time.Now(),
		Body:      []byte(str),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (q *RabbitMQ) Consume() (<-chan amqp.Delivery, error) {
	c, err := q.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return c, nil
}
func (q *RabbitMQ) Close() error {
	err := q.channel.Close()
	if err != nil {
		return err
	}
	return nil
}
