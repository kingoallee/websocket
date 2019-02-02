package main

import (
	"github.com/gorilla/websocket"
	"go-websocket/api"
	"net/http"
	"time"
)
var(
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandle(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("hello"))
	//io.WriteString(w, "nihao")
	var (
		wsConn *websocket.Conn
		err error
		data []byte
		conn *api.Connection
	)
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return
	}

	if conn, err = api.InitConnection(wsConn); err != nil {
		goto ERR
	}

	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()


	for {
		if  data, err = conn.ReadMessage();err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}

func main() {
	http.HandleFunc("/ws", wsHandle)
	http.ListenAndServe("localhost:9999", nil)
}
