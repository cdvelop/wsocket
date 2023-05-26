package wsocket_test

import (
	"time"

	"github.com/cdvelop/model"
	"github.com/cdvelop/wsocket"
	"github.com/gorilla/websocket"
)

func wsReply(hub *wsocket.WebSocket, ws *websocket.Conn) ([]model.Response, error) {
	var out []model.Response

	ws.SetReadDeadline(time.Now().Add(time.Millisecond * 10))

	_, msg, err := ws.ReadMessage()
	if err != nil {
		// log.Println("ReadMessage wsReply: ", err, msg)
		return out, err
	}

	return hub.DecodeResponses(msg), nil
}
