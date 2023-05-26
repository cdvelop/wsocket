package wsocket_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/cdvelop/model"
	"github.com/gorilla/websocket"
)

func newConn(a *model.User, token, origin string, server *httptest.Server) *websocket.Conn {

	url := url.URL{Scheme: "ws", Host: server.Listener.Addr().String(), Path: "/"}
	header := http.Header{"Authorization": []string{token}}
	header.Set("Origin", origin) // Establecer el origen en la cabecera
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), header)
	if err != nil {
		// log.Printf("Error al conectar al cliente %v %v", a.Name, err)
	}

	return conn

}
