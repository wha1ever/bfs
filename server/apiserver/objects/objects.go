package objects

import (
	"fmt"
	"github.com/wha1ever/bfs/server/apiserver/heartbeat"
	"io"
	"log"
	"net/http"
	"strings"
)

type PutStream struct {
	writer *io.PipeWriter
	c      chan error
}

func NewPutStream(server, object string) *PutStream {
	reader, writer := io.Pipe()
	c := make(chan error)
	go func() {
		request, _ := http.NewRequest("PUT", "http://"+server+"object/"+object, reader)
		client := http.Client{}
		r, err := client.Do(request)
		if err == nil && r.StatusCode != http.StatusOK {
			err = fmt.Errorf("dataServer return http code %d", r.StatusCode)
		}
		c <- err
	}()
	return &PutStream{
		writer: writer,
		c:      c,
	}
}
func (p *PutStream) Write(w []byte) (n int, err error) {
	return p.writer.Write(w)
}

func (p *PutStream) Close() error {
	err := p.writer.Close()
	if err != nil {
		return err
	}
	return <-p.c
}

func put(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	statusCode, err := storeObject(r.Body, object)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
}

func storeObject(r io.Reader, object string) (int, error) {
	stream, err := putStream(object)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	_, err = io.Copy(stream, r)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = stream.Close()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func putStream(object string) (*PutStream, error) {
	server := heartbeat.ChooseRandomDataServer()
	if server == "" {
		return nil, fmt.Errorf("cannot find any available dataServer")
	}
	return NewPutStream(server, object), nil
}
