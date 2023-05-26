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
						Object:           "error",
						Module:           "",
						Message:          "",
						Data:             []map[string]string{},
						SkipMeInResponse: false,
						Recipients:       []string{},
					},
				},
			}

			data_type, new_data_in, err := wsConn.ReadMessage()
			if err != nil {
				// log.Printf("CIERRE DE CONEXIÓN IP: %v DETALLE: %v\n", ip_ws_request, err)
				// h.NotifyListeners(h.SessionChange("login_out", &a))
				h.CLOSED_CONNECTION <- a
				h.UserRemove(a)
				break
			}

			switch data_type {

			case 1: // MENSAJE JSON TEXTO
				// log.Printf("MENSAJE: TEXTO %d\n", data_type)

				rq.Packages = h.DecodeResponses(new_data_in)

				// fmt.Println(">>> MENSAJE EN SERVIDOR ", rq)

				h.REQUESTS_IN <- &rq

			case 2:
				// fmt.Println(">>> MENSAJE BINARIO")

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
