package go_wsocket_test

import (
	"log"
	"time"

	"github.com/cdvelop/model"
	"github.com/gorilla/websocket"
)

func wsReply(ws *websocket.Conn) ([]*model.Response, error) {
	var out []*model.Response

	ws.SetReadDeadline(time.Now().Add(time.Millisecond * 10))
	_, msg, err := ws.ReadMessage()
	if err != nil {
		return out, err
	}

	out, err = hub.DecodeResponses(msg)
	if err != nil {
		log.Fatal(err)
	}

	return out, nil
}
