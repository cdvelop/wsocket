package wsocket

import (
	"sync"

	"github.com/cdvelop/cutkey"
	"github.com/cdvelop/model"
	"github.com/gorilla/websocket"
)

// concentrador de solicitantes para websocket
type WebSocket struct {
	upgrader *websocket.Upgrader

	concurrency_limiter int // Limitar la concurrencia de conexiones simultáneas ej 100
	// el mutex para proteger escritura de mapa de solicitantes
	lockUsers sync.RWMutex

	// Conexiones registradas.con
	users map[string]*model.User

	//canal de entrada nuevas solicitudes.
	REQUESTS_IN chan *model.Request

	// canal salida de solicitudes
	REQUESTS_OUT chan *model.Request

	// canal que avisa el cierre de las conexiones
	CLOSED_CONNECTION chan *model.User

	*cutkey.Cut
}
