package go_wsocket_test

import (
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cdvelop/model"
	"github.com/cdvelop/wsocket"
)

// en este test no tiene que llegar un mensaje de error a los
// demás solicitantes
func TestRespuestaError(t *testing.T) {

	//...
	const origin = "http://127.0.0.1"
	// Crear un nuevo servidor WebSocket
	hub = wsocket.New(&objects, 1024, 1, origin)

	// creamos solicitante A nivel acceso 3
	A := model.User{Token: "TOKEN_A", Ip: "", Name: "Maria", Area: 'v', AccessLevel: 3, Packages: make(chan []model.Response), LastConnection: time.Time{}}

	// creamos solicitante B nivel acceso 3
	B := model.User{Token: "TOKEN_B", Ip: "", Name: "Julio", Area: 'v', AccessLevel: 3, Packages: make(chan []model.Response), LastConnection: time.Time{}}

	go errorHandlerPrivateMessage(hub.REQUESTS_IN, hub.REQUESTS_OUT)

	// agregamos los solicitantes a hub
	hub.UserAdd(&A, &B)

	// iniciamos el servidor
	server := httptest.NewServer(hub)
	defer server.Close()

	// Conectar al servidor con el requirente A
	CONN_A := newConn(&A, A.Token, origin, server)

	// Conectar al servidor con el requirente B
	CONN_B := newConn(&B, B.Token, origin, server)

	// enviar mensaje sin nada
	message := model.Request{
		User: &A,
		Packages: []model.Response{
			{
				Type:             "",
				Object:           "chat",
				Module:           "",
				Message:          "",
				Data:             []map[string]string{},
				SkipMeInResponse: false,
				Recipients:       []string{},
			},
		},
	}

	sendMessage(CONN_A, &message)

	// respuesta A
	replies_A, err := wsReply(CONN_A)
	for _, reply_A := range replies_A {

		if err != nil && reply_A.Type != "error" {
			log.Fatal("No llego mensaje al destinatario A ", err)
		}
	}

	// fmt.Println("reply_A:", reply_A)

	// respuesta B
	_, err = wsReply(CONN_B)
	if err == nil {
		log.Fatal("ERROR MENSAJE LLEGO A: ", B.Name, " Y NO REALIZA NINGUNA SOLICITUD")
	}

	// segundo intento
	message = model.Request{
		User: &A,
		Packages: []model.Response{
			{
				Type:             "ok",
				Object:           "chat",
				Module:           "ok",
				Message:          "",
				Data:             []map[string]string{{"data": "123"}},
				SkipMeInResponse: false,
				Recipients:       []string{},
			},
		},
	}

	sendMessage(CONN_A, &message)

	// respuesta A
	replies_A, err = wsReply(CONN_A)
	if err != nil {
		log.Fatal("No llego mensaje al destinatario A ", err)
	}
	for _, reply_A := range replies_A {

		if reply_A.Type != "error" {
			log.Fatal("reply_A:", reply_A)
		}
	}
	// respuesta B
	_, err = wsReply(CONN_B)
	if err == nil {
		log.Fatal("ERROR MENSAJE LLEGO A: ", B.Name, " Y NO REALIZA NINGUNA SOLICITUD")
	}

}

func errorHandlerPrivateMessage(in <-chan *model.Request, out chan<- *model.Request) {

	select {
	case rq := <-in:
		for i, newPkg := range rq.Packages {

			// respondemos como error mal escrito
			newPkg.Type = "Errors"
			// fmt.Println("PROCESANDO SOLICITUD: ", rq.Type)
			rq.Packages[i] = newPkg

			out <- rq
		}
	case <-time.After(5 * time.Second): // Timeout after 5 seconds
		log.Fatal("error mensaje tardo mas de 5 seg. timed out")
	}

}
