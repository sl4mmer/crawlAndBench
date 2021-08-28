package bench

import (
	"context"
	"net"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

var count int32 = 0

func TestRun(t *testing.T) {
	router := http.NewServeMux()
	router.HandleFunc("/", handleRequest)
	srv := http.Server{
		Handler: router,
	}

	listener, err := net.Listen("tcp", ":8887")
	if err!=nil{
		t.Fail()
		return
	}
	defer listener.Close()
	go func() {
		srv.Serve(listener)
	}()
	result := Run(context.Background(), "rest://localhost:8887", 50,10)
	assert.Equal(t, result, 50)
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	v := atomic.LoadInt32(&count)
	if v >= 49 {
		writer.WriteHeader(429)
		writer.Write([]byte("bad"))
		return
	}
	writer.WriteHeader(200)
	writer.Write([]byte("ok"))
	atomic.AddInt32(&count, 1)
}
