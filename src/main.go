package main

import (
	"math"
	"strings"
	"unsafe"
	//"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math/big"
	"net/http"
	_ "strconv"

	"github.com/gorilla/websocket"
)

func findDelim(buffer []byte, idx int) int { // byte?? 区切り文字の位置を返すfunc
	delim := 0
	for i := idx; i < len(buffer); i++ {
		fmt.Printf("buffer[i] : %v\n", buffer[i])
		if 0 == buffer[i] {
			delim = i
			break
		}
	}
	//return byte(delim)
	fmt.Printf("Delim : %d\n", delim)
	return delim
}

// str, err := datas.buff.ReadString(nullstr)
type BufferReader struct {
	index int // current byte pos of delim
	buff  bytes.Buffer
	bytes []byte
}

func (b *BufferReader) readBigUInt64LE() (*big.Int, error) {
	lo, err := b.readUInt32LE()
	hi, err := b.readUInt32LE()
	if err != nil {
		return nil, err
	}
	biglow := big.NewInt(int64(lo))
	bighigh := big.NewInt(int64(hi))
	n := biglow.Or(biglow, bighigh)
	return n.Lsh(n, 32), nil
}

func (b *BufferReader) readUInt32LE() (uint32, error) {
	var i uint32
	buf := bytes.NewReader(b.buff.Bytes()[b.index:])
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		fmt.Println("binary.Read failed:", err)
		return 0, err
	}
	b.index += 4
	return i, nil
}

func (b *BufferReader) readInt32LE() int32 {
	val, err := b.buff.ReadBytes(byte(b.index))
	if err != nil {
		return -1
	}
	b.index += 4
	return *(*int32)(unsafe.Pointer(&val))
}

func (b *BufferReader) readInt16LE() int16 {
	val, err := b.buff.ReadBytes(byte(b.index))
	if err != nil {
		return -1
	}
	b.index += 2
	return *(*int16)(unsafe.Pointer(&val))
}

func (b *BufferReader) readInt8() int8 {
	val, err := b.buff.ReadBytes(byte(b.index))
	if err != nil {
		return -1
	}
	b.index++
	return *(*int8)(unsafe.Pointer(&val))
}

func (b *BufferReader) readUInt8() uint8 {
	val, err := b.buff.ReadBytes(byte(b.index))
	if err != nil {
		return 0 // todo
	}
	b.index++

	return *(*uint8)(unsafe.Pointer(&val))
}

func (b *BufferReader) readBoolean() bool {
	return 0 != b.readUInt8()
}

func (b *BufferReader) readFloatLE() (float32, error) {
	bits := b.bytes[b.index : b.index+4]
	uint32le := binary.LittleEndian.Uint32(bits)
	float := math.Float32frombits(uint32le)
	b.index += 4
	f := *(*float32)(unsafe.Pointer(&float))
	return f, nil
}

func (b *BufferReader) readCString() (string, error) {
	delim := findDelim(b.buff.Bytes(), b.index)
	var result string
	for b.index < delim {
		str, err := b.buff.ReadBytes(b.buff.Bytes()[delim])
		if err != nil {
			return "", err
		}
		b.index = delim + 1
		result = *(*string)(unsafe.Pointer(&str))
		result = strings.Trim(result, string(nullstr))
		//fmt.Printf("b.bytes : %v\nindex : %d\n", b.buff.Bytes(), b.index)
		return result, nil
	}
	//b.readCString()
	return "", nil
}

func (b *BufferReader) eof() bool {
	fmt.Printf("\nb.index : %d\nb.bytes len : %d\n", b.index, b.buff.Len())
	if b.index >= b.buff.Len() {
		fmt.Println("EOF")
		return true
	}
	return false

}

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

var (
	nullstr byte = []byte("\x00")[0]
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
	//defer websocket.Close()

	clients[websocket] = true

	SendRCON(websocket, "echo HELLO FROM GOLANG!")

	for {
		// メッセージ読み込み
		_, data, err := websocket.ReadMessage()
		if err != nil {
			log.Printf("error occurred while reading message: %v", err)
			delete(clients, websocket)
		}
		var datas = &BufferReader{
			index: 0,
			bytes: data,
		}
		datas.buff.Write(data)
		if !datas.eof() {
			cmd, err := datas.readCString()
			if err != nil {
				panic(err)
			}
			/*if datas.eof() {
				return
			}*/
			fmt.Printf("CMD : [%s]\n", cmd) //
			switch cmd {
			case "hello":
				//
			case "version":
				version, err := datas.readUInt32LE()
				if err != nil {
					panic(err)
				}
				fmt.Printf("VERSION : %d", version) //
			case "dataStop":
				//
			case "dataStart":
				//
			case "levelInit":
				mapname, err := datas.readCString()
				if err != nil {
					panic(err)
				}
				fmt.Printf("map : %s", mapname) //
				//
			case "levelShutdown":
				//
			case "cam":
				time, err := datas.readFloatLE()
				fmt.Printf("time = %f\n", time)
				xPosition, err := datas.readFloatLE()
				fmt.Printf("xPosition = %f\n", xPosition)
				yPosition, err := datas.readFloatLE()
				fmt.Printf("yPosition = %f\n", yPosition)
				zPosition, err := datas.readFloatLE()
				fmt.Printf("zPosition = %f\n", zPosition)
				xRotation, err := datas.readFloatLE()
				fmt.Printf("xRotation = %f\n", xRotation)
				yRotation, err := datas.readFloatLE()
				fmt.Printf("yRotation = %f\n", yRotation)
				zRotation, err := datas.readFloatLE()
				fmt.Printf("zRotation = %f\n", zRotation)
				fov, err := datas.readFloatLE()
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

func SendRCON(ws *websocket.Conn, cmd string) {
	command := []byte("exec")
	command = append(command, nullstr)
	command = append(command, []byte(cmd)...)
	command = append(command, nullstr)
	ws.WriteMessage(2, []uint8(command))
}
