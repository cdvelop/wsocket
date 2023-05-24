package go_wsocket_test

import (
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cdvelop/model"
	"github.com/cdvelop/wsocket"
)

var hub *wsocket.WebSocket

// Prueba unitaria de simulación de un chat básico privado entre 2 personas
// de una misma area y un tercero que no tiene que llegarle el mensaje
func TestPrivateMessage(t *testing.T) {
	const origin = "http://127.0.0.1"
	// Crear un nuevo servidor WebSocket
	hub = wsocket.New(&objects, 1024, 1, origin)

	//....

	// creamos solicitante area A
	A := model.User{Token: "TOKEN_A", Ip: "", Name: "Maria", Area: 'a', AccessLevel: 2, Packages: make(chan []*model.Response), LastConnection: time.Time{}}

	// creamos solicitante area A
	B := model.User{Token: "TOKEN_B", Ip: "", Name: "Julio", Area: 'a', AccessLevel: 3, Packages: make(chan []*model.Response), LastConnection: time.Time{}}

	// creamos intruso area C
	C := model.User{Token: "TOKEN_C", Ip: "", Name: "Espina", Area: 'c', AccessLevel: 4, Packages: make(chan []*model.Response), LastConnection: time.Time{}}

	go chatHandlerPrivateMessage(hub.REQUESTS_IN, hub.REQUESTS_OUT)

	// agregamos los solicitantes a hub
	hub.UserAdd(&A, &B, &C)

	// iniciamos el servidor
	server := httptest.NewServer(hub)
	defer server.Close()

	// Conectar al servidor con el requirente A
	CONN_A := newConn(&A, A.Token, origin, server)

	// Conectar al servidor con el requirente B
	CONN_B := newConn(&B, B.Token, origin, server)

	// Conectar al servidor con el intruso C
	CONN_C := newConn(&C, C.Token, origin, server)

	// Enviar un mensaje de "hola Maria" desde el requirente B al requirente A id secreto chat 111
	message := model.Request{
		User: &B,
		Packages: []*model.Response{
			{
				Type:    "create",
				Object:  "chat",
				Module:  "chat",
				Message: "",
				Data: []map[string]string{
					{"message": "hola Maria", "destination": "111"},
				},
				SkipMeInResponse: false,
				Recipients:       []string{},
			},
		},
	}

	sendMessage(CONN_B, &message)

	// respuesta A
	replies_A, err := wsReply(CONN_A)
	if err != nil {
		log.Fatal("No llego mensaje al destinatario A ", err)
	}
	for _, reply_A := range replies_A {

		if reply_A.Message != "" {
			log.Fatal("error llego mensaje confirmación a reply_A:", reply_A.Message)
		}
	}

	// respuesta B
	replies_B, err := wsReply(CONN_B)
	if err != nil {
		log.Fatal(err)
	}
	for _, reply_B := range replies_B {

		if reply_B.Message != expected_private_msg {
			log.Fatal("No me llego mensaje de ok")
		}
	}

	// fmt.Println("reply_B:", reply_B)

	// respuesta C
	_, err = wsReply(CONN_C)
	if err == nil {
		log.Fatal("ERROR MENSAJE LLEGO A: ", C.Name, " Y NO TIENE EL MISMO NIVEL")
	}

}

const expected_private_msg = "mensaje enviado ok"

func chatHandlerPrivateMessage(in <-chan *model.Request, out chan<- *model.Request) {

	// ej concentrador de llaves secretas del module
	room := map[string]string{
		"111": "TOKEN_A",
	}

	select {
	case rq := <-in:

		for _, pkg := range rq.Packages {

			switch pkg.Type {
			case "create":
				// fmt.Println("PROCESANDO SOLICITUD: ", pkg.Type)

				// **** 1-
				// añadimos data x a la solicitud u hacemos otro proceso
				pkg.Data = append(pkg.Data, map[string]string{"more_data": "xxxdata"})

				dest := room[pkg.Data[0]["destination"]]

				// log.Println("DESTINO: ", dest)
				// añadimos destinatario
				pkg.Recipients = append(pkg.Recipients, dest)

				// quitamos al emisor en la respuesta
				pkg.SkipMeInResponse = true

				//no es necesario agregar la respuesta si solo se envía a un destinatario ??
				// pkg.Packages = append(pkg.Packages, &pkg.Response)

				// la retornamos al destinatario
				out <- rq

				// *** 2-
				// creamos un nuevo rq para responder a emisor del mensaje
				resp := model.Request{
					User: rq.User,
					Packages: []*model.Response{{
						Type:             pkg.Type,
						Object:           pkg.Object,
						Module:           pkg.Module,
						Message:          expected_private_msg,
						Data:             []map[string]string{},
						SkipMeInResponse: false,
						Recipients:       []string{rq.Token},
					}},
				}

				out <- &resp

			}

		}

	case <-time.After(5 * time.Second): // Timeout after 5 seconds
		log.Fatal("error mensaje tardo mas de 5 seg. timed out")
	}

}
