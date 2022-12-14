package main

import (
	"context"
	"github.com/wha1ever/bfs/server/apiserver/heartbeat"
	"github.com/wha1ever/bfs/server/apiserver/locate"
	"github.com/wha1ever/bfs/server/apiserver/object"
	"log"
	"net/http"
	"os"
)

func main() {

	go func() {
		err := heartbeat.ListenHeartbeat(context.Background())
		if err != nil {
			log.Println(err)
		}
	}()
	http.HandleFunc("/object/", object.Handler)
	http.HandleFunc("/locate/", locate.Handler)

	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
