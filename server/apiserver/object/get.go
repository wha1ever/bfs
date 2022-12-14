package object

import (
	"fmt"
	"github.com/wha1ever/bfs/server/apiserver/locate"
	"golang.org/x/net/context"
	"io"
	"log"
	"net/http"
	"strings"
)

type GetStream struct {
	reader io.Reader
}

func (g *GetStream) Read(p []byte) (n int, err error) {
	//TODO implement me
	//panic("implement me")
	return g.reader.Read(p)
}

func NewGetStream(server, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}

	return newGetStream("http://" + server + "/objects/" + object)
}

func newGetStream(url string) (*GetStream, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	return &GetStream{reader: r.Body}, nil
}

func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := getStream(r.Context(), object)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = io.Copy(w, stream)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func getStream(ctx context.Context, object string) (io.Reader, error) {
	server, err := locate.Locate(ctx, object)
	if err != nil {
		return nil, err
	}
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return NewGetStream(server, object)
}