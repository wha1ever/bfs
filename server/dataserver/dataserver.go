package main

import (
	"context"
	"github.com/wha1ever/bfs/server/dataserver/heartbeat"
	"github.com/wha1ever/bfs/server/dataserver/locate"
	"github.com/wha1ever/bfs/server/dataserver/object"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	go func() {
		heartbeat.StartHeartbeat()
	}()
	go func() {
		err := locate.StartLocate(ctx)
		if err != nil {
			log.Println(err)
		}
	}()
	http.HandleFunc("/object/", object.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
