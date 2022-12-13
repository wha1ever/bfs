package main

import (
	"github.com/wha1ever/bfs/server/dataserver/heartbeat"
	"github.com/wha1ever/bfs/server/dataserver/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.StartHeartbeat()

	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
