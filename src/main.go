package main

import (
	"Apush"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

func main() {
	data, queryStr := Apush.GetQueryString()

	path := "ws://10.125.198.202:6080/apush/1/websocket/" + data[0] + "," + data[6] + queryStr
	fmt.Println(path)

	c, _, err := websocket.DefaultDialer.Dial(path, nil)
	if err != nil {
		fmt.Println("[ERROR] ", err)
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
		fmt.Println(string(message))
		if strings.HasPrefix(string(message), "2::") {
			c.WriteMessage(websocket.TextMessage, []byte("2::"))
		}
	}
	// 需要实现Client和Server之间的协议
	// apush.client.external.io.socket.IOConnection#transportMessage
	c.Close()
}
