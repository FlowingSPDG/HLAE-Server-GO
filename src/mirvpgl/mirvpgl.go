package mirvpgl

import (
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

// サーバー起動
func main() {
	http.Handle("/echo", websocket.Handler(EchoHandler))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func EchoHandler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}
