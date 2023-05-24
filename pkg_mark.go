package wsocket

import (
	"strings"

	"github.com/cdvelop/model"
)

func markPackages(rq *model.Request) (pkg_error map[int]*model.Response, pkg_public map[int]*model.Response, pkg_area map[int]*model.Response, pkg_private map[int]*model.Response) {

	pkg_error = make(map[int]*model.Response, len(rq.Packages))
	pkg_private = make(map[int]*model.Response, len(rq.Packages))
	pkg_area = make(map[int]*model.Response, len(rq.Packages))
	pkg_public = make(map[int]*model.Response, len(rq.Packages))
	// 1- marcación de paquetes según tipo: error, privado, area o publico
	for id, pkg := range rq.Packages {

		if strings.Contains(strings.ToLower(pkg.Type), "err") {
			// log.Println("marcación error destinatario el emisor")
			pkg.Type = "error"
			pkg_error[id] = pkg

		} else {

			if len(pkg.Recipients) != 0 {
				// log.Println("marcación privada destinatarios específicos")
				pkg_private[id] = pkg

			} else if pkg.TargetArea { //
				// log.Println("marcación por Area")
				pkg_area[id] = pkg
			} else {
				// log.Println("marcación pública")
				pkg_public[id] = pkg

			}
		}
	}

	return
}
