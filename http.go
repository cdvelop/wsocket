package wsocket

import (
	"log"
	"net/http"
	"sync"
)

func (h *WebSocket) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Intentar adquirir un valor del semáforo antes de procesar la solicitud
	h.concurrency_limiter <- struct{}{}
	defer func() { <-h.concurrency_limiter }()

	token := r.Header.Get("Authorization")
	// fmt.Println("TOKEN ENVIADO: ", token)

	user, exists := h.users[token]
	if !exists || user == nil {
		http.Error(w, "Acceso no autorizado", http.StatusUnauthorized)
		return
	}

	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade WebSocket Error %v %s\n", r.RemoteAddr, err)
		return
	}

	defer h.UserRemove(user)
	var wg sync.WaitGroup
	// WaitGroup espera a que termine una colección de gorutinas.
	// La gorutina principal llama a Agregar para establecer el número de gorutinas a esperar.
	// Luego, cada una de las gorutinas corre y llama Hecho cuando termina. Al mismo tiempo,
	// Un WaitGroup no debe copiarse después del primer uso.
	wg.Add(2)

	go h.listen(user, &wg, ws)
	go h.respond(user, &wg, ws)

	wg.Wait()
	// Wait puede usarse para bloquear hasta que todas las gorutinas hayan terminado.
	ws.Close()
}
