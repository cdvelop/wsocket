package wsocket_test

import (
	"github.com/cdvelop/model"
)

var objects = []model.Object{
	{
		Name: "user",
		Fields: []model.Field{
			{Name: "name"},
			{Name: "email"},
			{Name: "phone"},
		},
	},
	{
		Name: "chat",
		Fields: []model.Field{
			{Name: "message"},
			{Name: "destination"},
		},
	},
}
