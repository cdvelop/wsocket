package wsocket

import (
	"fmt"
	"log"
	"sync"

	"github.com/cdvelop/model"
	"github.com/gorilla/websocket"
)

func (h *WebSocket) listen(u *model.User, wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for {

		ip_ws_request := fmt.Sprintf("%v", wsConn.RemoteAddr())

		// log.Printf("SOLICITANTE: IP: %v\n", ip_ws_request)

		if a, exist := h.users[u.Token]; exist && a != nil {

			// log.Printf("NOMBRE: %v\n", u.Name)
			// si esta todo correcto asignamos el puntero del requirente a la solicitud
			var rq = model.Request{
				User: u,
				Packages: []model.Response{
					{
						Type:             "error", //inicializamos por defecto en caso de error
						Object:           "",
						Module:           "",
						Message:          "",
						Data:             []map[string]string{},
						SkipMeInResponse: false,
						Recipients:       []string{},
					},
				},
			}

			messageType, messageData, err := wsConn.ReadMessage()
			if err != nil {
				// log.Printf("CIERRE DE CONEXIÓN IP: %v DETALLE: %v\n", ip_ws_request, err)
				// h.NotifyListeners(h.SessionChange("login_out", &a))
				h.CLOSED_CONNECTION <- a
				h.UserRemove(a)
				break
			}

			switch messageType {

			case 1: // MENSAJE JSON TEXTO
				// log.Printf("MENSAJE: TEXTO %d\n", messageType)

				rq.Packages, err = h.DecodeResponses(messageData)
				if err != nil {
					rq.Packages[0].Message = err.Error()
					u.Packages <- rq.Packages
					break
				}

				if rq.Packages[0].Type == "" || rq.Packages[0].Object == "" || rq.Packages[0].Module == "" || len(rq.Packages[0].Data) == 0 {

					rq.Packages[0].Type = "error"
					rq.Packages[0].Message = "error solicitud sin información"
					// fmt.Println(rq.Message)
					u.Packages <- rq.Packages
					break
				}

				if rq.Token == u.Token {
					// fmt.Println(">>> MENSAJE EN SERVIDOR ", rq.Message)

					h.REQUESTS_IN <- &rq

				} else {
					// fmt.Println(">>> LLAVE CONEXIÓN DIFERENTE <<<")
					rq.Packages[0].Type = "error"
					rq.Packages[0].Message = "Error Acceso no autorizado"
					u.Packages <- rq.Packages
					wsConn.Close()
				}

			case 2: // MENSAJE BINARIO
				rq.Packages[0].Type = "error"
				rq.Packages[0].Message = "Error Archivos Binario no soportado"
				u.Packages <- rq.Packages
			}

		} else {

			wsConn.Close()
			log.Printf("\n¡ERROR INTENTO DE INGRESO CLIENTE %v NO REGISTRADO!", ip_ws_request)
			break
		}
	}
}
