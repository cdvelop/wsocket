package go_wsocket_test

import (
	"log"

	"github.com/cdvelop/model"
	"github.com/gorilla/websocket"
)

func sendMessage(ws *websocket.Conn, rq *model.Request) {
	// t.Helper()

	out, err := hub.EncodeResponses(rq.Packages)
	if err != nil {
		log.Fatal("!!!sendMessage EncodeResponses Error ", err)
	}

	if err := ws.WriteMessage(websocket.TextMessage, out); err != nil {
		log.Fatal("sendMessage", err)
	}
}
