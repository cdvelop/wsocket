package wsocket_test

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cdvelop/model"
	"github.com/cdvelop/wsocket"
	"github.com/gorilla/websocket"
)

func TestBinaryFile(t *testing.T) {
	// Crear un servidor WebSocket
	hub := wsocket.New(objects, 1024, 1000, origin)
	// creamos solicitante area A
	A := model.User{Token: "TOKEN_A", Ip: "", Name: "Maria", Area: 'a', AccessLevel: 2, Packages: make(chan []model.Response), LastConnection: time.Time{}}

	// agregamos los solicitantes a hub
	hub.UserAdd(&A)

	// iniciamos el servidor
	server := httptest.NewServer(hub)
	defer server.Close()

	// Conectar al servidor con el requirente A
	USER_A := newConn(&A, A.Token, origin, server)

	// enviamos el formulario como binario
	if err := USER_A.WriteMessage(websocket.BinaryMessage, binaryTestFile()); err != nil {
		log.Fatal("sendBinaryMessage", err)
	}

	// respuesta del envió del archivo
	REPLIES_USER_A, _ := wsReply(hub, USER_A)

	for i, REPLY_A := range REPLIES_USER_A {

		if i > 0 {
			log.Fatal("se esperaba solo un mensaje")
		}

		if REPLY_A.Message != "Error Archivos Binario no soportado" {
			log.Fatalln("expectativa no cumplida")
		}

	}

}

func binaryTestFile() []byte {

	const file_name = "03-dino-test.png"

	// Leer el archivo
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal("openFile:", err)
	}
	defer file.Close()

	// Obtener el tamaño del archivo
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal("fileStat:", err)
	}
	fileSize := fileInfo.Size()

	// Crear un slice de bytes para almacenar el contenido del archivo
	fileBytes := make([]byte, fileSize)

	// Leer el contenido del archivo en el slice de bytes
	_, err = file.Read(fileBytes)
	if err != nil {
		log.Fatal("readFile:", err)
	}

	return fileBytes
}
