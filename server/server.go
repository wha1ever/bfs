package main

import (
	"github.com/wha1ever/bfs/objects"
	"log"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
