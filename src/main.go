package main

import (
	"fmt"
	"github.com/FlowingSPDG/HLAE-Server-GO/HLAE-Server-GO/src/mirvpgl"
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
	// websocket の状態を更新
	websocket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error upgrading GET request to a websocket::", err)
	}
	// websocket を閉じる
	//defer websocket.Close()

	clients[websocket] = true

	err = mirvpgl.SendRCON(websocket, "echo HELLO FROM GOLANG!")
	if err != nil {
		panic(err)
	}

	for {
		// メッセージ読み込み
		datatype, data, err := websocket.ReadMessage()
		if err != nil {
			log.Printf("error occurred while reading message: %v", err)
			websocket.Close()
			delete(clients, websocket)
			break
		}
		if datatype == CloseMessage {
			log.Println("CloseMessage received")
			delete(clients, websocket)
			break
		}
		var datas = &mirvpgl.BufferReader{
			Index: 0,
			Bytes: data,
		}
		datas.Buff.Write(data)
		if !datas.Eof() {
			cmd, err := datas.ReadCString()
			if err != nil {
				panic(err)
			}
			fmt.Printf("CMD : [%s]\n", cmd) //
			switch cmd {
			case "hello":
				//
			case "version":
				version, err := datas.ReadUInt32LE()
				if err != nil {
					panic(err)
				}
				fmt.Printf("VERSION : %d", version) //
			case "dataStop":
				//
			case "dataStart":
				//
			case "levelInit":
				mapname, err := datas.ReadCString()
				if err != nil {
					panic(err)
				}
				fmt.Printf("map : %s", mapname) //
				//
			case "levelShutdown":
				//
			case "cam":
				time, err := datas.ReadFloatLE()
				fmt.Printf("time = %f\n", time)
				xPosition, err := datas.ReadFloatLE()
				fmt.Printf("xPosition = %f\n", xPosition)
				yPosition, err := datas.ReadFloatLE()
				fmt.Printf("yPosition = %f\n", yPosition)
				zPosition, err := datas.ReadFloatLE()
				fmt.Printf("zPosition = %f\n", zPosition)
				xRotation, err := datas.ReadFloatLE()
				fmt.Printf("xRotation = %f\n", xRotation)
				yRotation, err := datas.ReadFloatLE()
				fmt.Printf("yRotation = %f\n", yRotation)
				zRotation, err := datas.ReadFloatLE()
				fmt.Printf("zRotation = %f\n", zRotation)
				fov, err := datas.ReadFloatLE()
				fmt.Printf("fov = %f\n", fov)
				if err != nil {
					panic(err)
				}
				//
			case "gameEvent":
				//TODO. JSON
			default:
				fmt.Println("Unknown message")
			}
		}
	}
}

func main() {
	http.HandleFunc("/mirv", HandleClients)
	err := http.ListenAndServe(":63337", nil)
	if err != nil {
		log.Fatal("error starting http server::", err)
		return
	}
}
