package wsocket

import (
	"github.com/cdvelop/model"
)

func addPackage(user *model.User, rq *model.Request, id int, pkg *model.Response, pkg_reg *map[string]pkgReg) {
	var jump_token string

	if pkg.SkipMeInResponse {
		jump_token = rq.Token
	}
	if user.Token != jump_token {
		data, exist := (*pkg_reg)[user.Token]
		if !exist { // inicializamos las variables si no existe el usuario
			data.registered = make(map[int]struct{})
			data.out_pkgs = make([]*model.Response, 0)
		}

		if _, exist := data.registered[id]; !exist {

			prepareMessage(user, rq, pkg)

			data.out_pkgs = append(data.out_pkgs, pkg)
			data.registered[id] = struct{}{}

			(*pkg_reg)[user.Token] = data

		}

	}
}

func (h *WebSocket) addPrivatePackages(rq *model.Request, pkg_private *map[int]*model.Response, pkg_reg *map[string]pkgReg) {

	for id, pkg := range *pkg_private {

		for _, destination_address := range pkg.Recipients {
			if user, exists := h.users[destination_address]; exists {

				addPackage(user, rq, id, pkg, pkg_reg)

			}
		}
	}
}

func prepareMessage(u *model.User, rq *model.Request, pkg *model.Response) {
	message := pkg.Message

	if message != "" && u.Token == rq.Token {
		pkg.Message = message
	} else {
		pkg.Message = ""
	}
}
