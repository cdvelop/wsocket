package wsocket

import (
	"log"
	"sync"

	"github.com/cdvelop/model"
	"github.com/gorilla/websocket"
)

func (h *WebSocket) respond(a *model.User, wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for resp := range a.Packages {

		err := wsConn.WriteMessage(websocket.TextMessage, h.EncodeResponses(resp))
		if err != nil {
			log.Printf("!!! WriteMessage User: %v Error: %v\n", a.Name, err)
			break
		}
	}
}
