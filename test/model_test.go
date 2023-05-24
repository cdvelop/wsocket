package go_wsocket_test

import (
	"github.com/cdvelop/model"
)

var objects = map[string]model.Object{
	"user": {
		Name: "Usuario",
		Fields: []model.Field{
			{Name: "name"},
			{Name: "email"},
			{Name: "phone"},
		},
	},
	"chat": {
		Name: "chat",
		Fields: []model.Field{
			{Name: "message"},
			{Name: "destination"},
		},
	},
}
