package locate

import (
	"context"
	"encoding/json"
	"github.com/wha1ever/bfs/internal/rabbitmq"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Locate(ctx context.Context, name string) (string, error) {
	q := rabbitmq.New(ctx, os.Getenv("RABBITMQ_SERVER"))
	err := q.Publish("dataServers", name)
	if err != nil {
		return "", err
	}
	c, err := q.Consume()
	if err != nil {
		log.Println(err)
		return "", err
	}
	go func() {
		time.Sleep(time.Second)
		err := q.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s, nil

}
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info, err := Locate(r.Context(), strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(info)
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Exist(ctx context.Context, name string) bool {
	s, err := Locate(ctx, name)
	if err != nil {
		log.Println(err)
		return false
	}
	return s != ""
}
