package wsocket_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cdvelop/model"
	"github.com/cdvelop/wsocket"
)

func TestStress(t *testing.T) {

	concurrencia := 100
	// Crear un servidor WebSocket
	var hub2 = wsocket.New(objects, 1024, concurrencia, origin)

	// iniciamos el servidor
	server := httptest.NewServer(hub2)
	defer server.Close()

	A := model.User{Token: "TOKEN_A", Ip: "", Name: "Maria", Area: 'a', AccessLevel: 2, Packages: make(chan []model.Response), LastConnection: time.Time{}}

	for i := 0; i < concurrencia; i++ {

		// go func() {

		newConn(&A, A.Token, origin, server)

		// }()

	}

}
