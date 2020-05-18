package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/googollee/go-socket.io"
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		fmt.Println("Cannnot create server instance.")
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected ID: ", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/", "post", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		fmt.Println(msg)
		return "recv" + msg
	})

	server.OnError("/", func(s socketio.Conn, err error) {
		fmt.Println("Internal Server Error.", err)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	server.BroadcastToRoom("", "bcast", "allPosts", "hoge")

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	log.Fatal(http.ListenAndServe(":8080", nil))
}