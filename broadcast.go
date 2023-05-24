package wsocket

import (
	"fmt"
	"time"

	"github.com/cdvelop/model"
)

type pkgReg struct {
	registered map[int]struct{}
	out_pkgs   []*model.Response
}

func (h *WebSocket) StartBroadcasting(in <-chan *model.Request) {
	//respuesta publica
	go func() {

		for rq := range in {

			h.lockUsers.Lock()

			pkg_error, pkg_public, pkg_area, pkg_private := markPackages(rq)

			// si solo hay errores retornamos todo al emisor y finalizamos la solicitud
			if len(pkg_error) == len(rq.Packages) {
				rq.User.Packages <- rq.Packages
			} else {

				// 2- distribución

				// Mapa para recolectar los id de los paquetes que se deben enviar a cada usuario
				pkg_reg := map[string]pkgReg{}

				if len(pkg_public) != 0 || len(pkg_area) != 0 { // 2.1 pública o por area

					// recorremos todos los usuarios y ademas
					// agregaremos los paquetes area o privados si son diferentes
					for _, user := range h.users {

						for id, pkg := range pkg_public {
							addPackage(user, rq, id, pkg, &pkg_reg)
						}

						for id, pkg := range pkg_area {
							addPackage(user, rq, id, pkg, &pkg_reg)
						}

					}

					h.addPrivatePackages(rq, &pkg_private, &pkg_reg)

				} else { // 2.3 privado buscamos solo los usuarios específicos

					h.addPrivatePackages(rq, &pkg_private, &pkg_reg)

				}

				// hay algún paquete con error pendiente para adjuntarlo al usuario emisor?
				for _, pkg := range pkg_error {
					data, exist := pkg_reg[rq.User.Token]
					if !exist {
						data.registered = make(map[int]struct{})
						data.out_pkgs = make([]*model.Response, 0)
					}
					data.out_pkgs = append(data.out_pkgs, pkg)
					pkg_reg[rq.User.Token] = data
				}

				// envió de información
				for address, data := range pkg_reg {
					user, exist := h.users[address]
					if exist {
						h.respondRequest(user, data.out_pkgs)
					}
				}

			}

			h.lockUsers.Unlock()

		}
	}()
}

func (h *WebSocket) respondRequest(user *model.User, responses []*model.Response) {

	select {
	case user.Packages <- responses:
		// deja de intentar enviar a esta conexión después de intentarlo durante 1 segundo.
		// si tenemos que parar, significa que un lector murió así que quita la conexión también.
	case <-time.After(2 * time.Second):
		fmt.Printf("Time After cierre de la conexión %v\n", user.Ip)
		h.UserRemove(user)
		// default:
		// 	fmt.Printf("Default cierre de la conexión %v\n", user.credentials)
		// 	user.hub.userRemove(user)
	}
}
