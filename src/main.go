package main

import (
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

func findDelim(buffer []byte, idx int) int { // byte??
	delim := 0
	for i := idx; i < len(buffer); i++ {
		fmt.Printf("buffer[i] : %v\n", buffer[i])
		if nullstr == buffer[i] {
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
	index int
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
	buf := bytes.NewReader(b.bytes[b.index:])
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		fmt.Println("binary.Read failed:", err)
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

func (b *BufferReader) readFloatLE() (float64, error) {
	val, err := b.buff.ReadBytes(byte(b.index))
	if err != nil {
		return -1, err
	}
	return *(*float64)(unsafe.Pointer(&val)), nil
}

func (b *BufferReader) readCString() ([]string, error) {
	delim := findDelim(b.bytes, b.index)
	var result []string
	for i := delim; b.index <= delim; i++ {
		str, err := b.buff.ReadBytes(byte(b.index))
		if err != nil {
			return nil, err
		}
		result = append(result, *(*string)(unsafe.Pointer(&str)))
		/*if err != nil {
			return "", err
		}*/
		b.index = delim + 1
		fmt.Printf("b.bytes : %v\nindex : %d\n", b.bytes, b.index)
	}
	return result, nil
}

func (b *BufferReader) eof() bool {
	return b.index >= len(b.bytes)
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
	defer websocket.Close()

	clients[websocket] = true

	command := []byte("exec")
	command = append(command, nullstr)
	command = append(command, []byte("echo hello from GOLANG")...)
	command = append(command, nullstr)
	websocket.WriteMessage(2, []uint8(command))

	// メッセージ読み込み
	_, data, err := websocket.ReadMessage()
	if err != nil {
		log.Printf("error occurred while reading message: %v", err)
		delete(clients, websocket)
	}
	var datas = BufferReader{
		index: 0,
		bytes: data,
	}
	datas.buff.Write(data)
	for !datas.eof() {
		cmd, err := datas.readCString()
		ver, err := datas.readUInt32LE()
		if err != nil {
			panic(err)
		}

		//str, err := datas.buff.ReadString(nullstr)
		//str, _ := datas.readCString()
		fmt.Printf("CMD : %s", cmd) //
		fmt.Printf("Version : %d", ver)

		//datastr := string(data)
		//fmt.Println(datastr)
		//fmt.Println(datatype)
	}
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
