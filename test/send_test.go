package wsocket_test

import (
	"log"

	"github.com/cdvelop/model"
	"github.com/cdvelop/wsocket"
	"github.com/gorilla/websocket"
)

func sendMessage(hub *wsocket.WebSocket, ws *websocket.Conn, rq *model.Request) {

	data := hub.EncodeResponses(rq.Packages)

	if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println("sendMessage", err)
	}
}
