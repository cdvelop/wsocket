package wsocket

import "github.com/cdvelop/model"

func (h *WebSocket) UserAdd(users ...*model.User) {

	h.lockUsers.Lock()
	defer h.lockUsers.Unlock()

	for _, a := range users {
		h.users[a.Token] = a
	}

}

func (hub *WebSocket) UserRemove(a *model.User) {
	hub.lockUsers.Lock()
	if user, ok := hub.users[a.Token]; ok {
		delete(hub.users, user.Token)
		close(user.Packages)
	}
	defer hub.lockUsers.Unlock()
}
