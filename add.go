package wsocket

import (
	"net/http"
	"sync"

	"github.com/cdvelop/cutkey"
	"github.com/cdvelop/model"
	"github.com/gorilla/websocket"
)

// models mapa con los modelos de objetos que utiliza el sistema
// allowed_origins ej: "http://localhost", "http://127.0.0.1", "http://example.com", "https://example.com"
// buffer_size ej: 1024
// concurrency_max Limitar la concurrencia de conexiones simultáneas ej 100
func New(models *[]model.Object, buffer_size, concurrency_max int, allowed_origins ...string) *WebSocket {

	ws := WebSocket{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  buffer_size,
			WriteBufferSize: buffer_size,
			CheckOrigin: func(r *http.Request) bool {
				allowedOrigins := allowed_origins
				origin := r.Header.Get("Origin")
				// fmt.Println("ORIGIN: ", origin)
				for _, allowed := range allowedOrigins {
					if allowed == origin {
						return true
					}
				}
				return false
			},
		},
		concurrency_limiter: concurrency_max,
		lockUsers:           sync.RWMutex{},
		users:               map[string]*model.User{},
		REQUESTS_IN:         make(chan *model.Request),
		REQUESTS_OUT:        make(chan *model.Request),
		CLOSED_CONNECTION:   make(chan *model.User),

		Add: cutkey.Add{
			Objects: models,
		},
	}

	ws.StartBroadcasting(ws.REQUESTS_OUT)

	return &ws
}
