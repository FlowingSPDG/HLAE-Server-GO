package main

import (
	"fmt"
	"log"
	"net/http"
	_ "strconv"

	"github.com/gorilla/websocket"
)

// The message types are defined in RFC 6455, section 11.8.
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)

// WebSocket サーバーにつなぎにいくクライアント
var clients = make(map[*websocket.Conn]bool)

// WebSocket 更新用
var upgrader = websocket.Upgrader{}

// クライアントのハンドラ
func HandleClients(w http.ResponseWriter, r *http.Request) {
	// ゴルーチンで起動
	//broadcastMessagesToClients()
	// websocket の状態を更新
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error upgrading GET request to a websocket::", err)
	}
	// websocket を閉じる
	defer websocket.Close()

	clients[websocket] = true

	command := []byte("exec")
	command = append(command, []byte("\x00")...)
	command = append(command, []byte("echo hello from GOLANG")...)
	command = append(command, []byte("\x00")...)
	newStr := []uint8(command)
	websocket.WriteMessage(2, newStr)
	//ws.send(new Uint8Array(Buffer.from('exec\0'+data.trim()+'\0','utf8')),{binary: true});

	// メッセージ読み込み
	datatype, data, err := websocket.ReadMessage()
	if err != nil {
		log.Printf("error occurred while reading message: %v", err)
		delete(clients, websocket)
	}
	//datauint8 := []uint8(data)
	datastr := string(data)
	fmt.Println(datastr)
	fmt.Println(datatype)
	fmt.Println(data)

}

func main() {
	// localhost:8080 でアクセスした時に index.html を読み込む

	http.HandleFunc("/mirv", HandleClients)
	err := http.ListenAndServe(":63337", nil)
	if err != nil {
		log.Fatal("error starting http server::", err)
		return
	}
}
