package mirvpgl

import (
	"math"
	"strings"
	"unsafe"
	//"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/big"
	"net/http"
	// "strconv"
)

var (
	nullstr []byte = []byte("\x00")

	// WebSocket 更新用
	upgrader = websocket.Upgrader{}

	// WebSocket サーバーにつなぎにいくクライアント
	clients = make(map[*websocket.Conn]bool)
)

func FindDelim(buffer []byte, idx int) int { // byte?? 区切り文字の位置を返すfunc
	delim := 0
	for i := idx; i < len(buffer); i++ {
		// fmt.Printf("buffer[i] : %v\n", buffer[i])
		if 0 == buffer[i] {
			delim = i
			break
		}
	}
	// return byte(delim)
	// fmt.Printf("Delim : %d\n", delim)
	return delim
}

// str, err := datas.buff.ReadString(nullstr)
type BufferReader struct {
	Index int // current byte pos of delim
	Buff  bytes.Buffer
	Bytes []byte
}

func (b *BufferReader) ReadBigUInt64LE() (*big.Int, error) {
	lo, err := b.ReadUInt32LE()
	hi, err := b.ReadUInt32LE()
	if err != nil {
		return nil, err
	}
	biglow := big.NewInt(int64(lo))
	bighigh := big.NewInt(int64(hi))
	n := biglow.Or(biglow, bighigh)
	return n.Lsh(n, 32), nil
}

func (b *BufferReader) ReadUInt32LE() (uint32, error) {

	bits := b.Bytes[b.Index : b.Index+4]
	uint32le := binary.LittleEndian.Uint32(bits)
	b.Index += 4
	f := *(*uint32)(unsafe.Pointer(&uint32le))
	return f, nil

	var i uint32
	buf := bytes.NewReader(b.Buff.Bytes()[b.Index:])
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		fmt.Println("binary.Read failed:", err)
		return 0, err
	}
	b.Index += 4
	return i, nil
}

func (b *BufferReader) ReadInt32LE() int32 {
	val, err := b.Buff.ReadBytes(byte(b.Index))
	if err != nil {
		return -1
	}
	b.Index += 4
	return *(*int32)(unsafe.Pointer(&val))
}

func (b *BufferReader) ReadInt16LE() int16 {
	val, err := b.Buff.ReadBytes(byte(b.Index))
	if err != nil {
		return -1
	}
	b.Index += 2
	return *(*int16)(unsafe.Pointer(&val))
}

func (b *BufferReader) ReadInt8() int8 {
	val, err := b.Buff.ReadBytes(byte(b.Index))
	if err != nil {
		return -1
	}
	b.Index++
	return *(*int8)(unsafe.Pointer(&val))
}

func (b *BufferReader) ReadUInt8() uint8 {
	val := b.Bytes[b.Index : b.Index+1]
	b.Index++
	return *(*uint8)(unsafe.Pointer(&val))
}

func (b *BufferReader) ReadBoolean() bool {
	return 0 != b.ReadUInt8()
}

func (b *BufferReader) ReadFloatLE() (float32, error) {
	bits := b.Bytes[b.Index : b.Index+4]
	uint32le := binary.LittleEndian.Uint32(bits)
	float := math.Float32frombits(uint32le)
	b.Index += 4
	f := *(*float32)(unsafe.Pointer(&float))
	return f, nil
}

func (b *BufferReader) ReadCString() (string, error) {
	delim := FindDelim(b.Buff.Bytes(), b.Index)
	var result string
	for b.Index < delim {
		str, err := b.Buff.ReadBytes(b.Buff.Bytes()[delim])
		if err != nil {
			return "", err
		}
		b.Index = delim + 1
		result = *(*string)(unsafe.Pointer(&str))
		result = strings.Trim(result, string(nullstr))
		//fmt.Printf("b.bytes : %v\nindex : %d\n", b.buff.Bytes(), b.index)
		return result, nil
	}
	return "", nil
}

func (b *BufferReader) Eof() bool {
	// fmt.Printf("\nb.index : %d\nb.bytes len : %d\n", b.Index, len(b.Bytes))
	if b.Index >= b.Buff.Len() {
		// fmt.Println("EOF")
		return true
	}
	return false
}

type CamData struct {
	Time float32
	Fov  float32
	XPos float32
	YPos float32
	ZPos float32
	XRot float32
	YRot float32
	Zrot float32
}

type HLAEServer struct {
	ws          websocket.Conn
	handlers    []func(string)
	camhandlers []func(*CamData)
}

// SendRCON command
func (h *HLAEServer) SendRCON(cmd string) error {
	length := len("exec") + len(nullstr) + len(cmd) + len(nullstr)
	command := make([]byte, 0, length)
	command = append(command, []byte("exec")...)
	command = append(command, nullstr...)
	command = append(command, []byte(cmd)...)
	command = append(command, nullstr...)
	err := h.ws.WriteMessage(2, []uint8(command))
	if err != nil {
		return err
	}
	return nil
}

// RegisterHandler to handle each requests
func (h *HLAEServer) RegisterHandler(handler func(string)) {
	if h.handlers == nil {
		h.handlers = make([]func(string), 0)
	}
	h.handlers = append(h.handlers, handler)
	log.Printf("Registered handler. Currently %d handlers are active\n", len(h.handlers))
}

// RegisterCamHandler to handle each requests
func (h *HLAEServer) RegisterCamHandler(handler func(*CamData)) {
	if h.camhandlers == nil {
		h.camhandlers = make([]func(*CamData), 0)
	}
	h.camhandlers = append(h.camhandlers, handler)
	log.Printf("Registered Camera handler. Currently %d handlers are active\n", len(h.handlers))
}

func (h *HLAEServer) handleRequest(cmd string) {
	log.Printf("Sending handler request for %d clients...\n", len(h.handlers))
	for i := 0; i < len(h.handlers); i++ {
		go h.handlers[i](cmd)
	}
}

func (h *HLAEServer) handleCamRequest(cam *CamData) {
	log.Printf("Sending camera handler request for %d clients...\n", len(h.handlers))
	for i := 0; i < len(h.handlers); i++ {
		go h.camhandlers[i](cam)
	}
}

// Start listening
func (h *HLAEServer) Start(host, path string) {
	log.Printf("Listening on %s%s", host, path)
	http.HandleFunc(path, h.handleClients)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		panic(err)
	}
}

// クライアントのハンドラ
func (h *HLAEServer) handleClients(w http.ResponseWriter, r *http.Request) {
	// websocket の状態を更新
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error upgrading GET request to a websocket::", err)
	}
	// websocket を閉じる
	defer ws.Close()

	clients[ws] = true
	h.ws = *ws

	err = h.SendRCON("echo HELLO FROM GOLANG!")
	if err != nil {
		panic(err)
	}

	for {
		// メッセージ読み込み
		datatype, data, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error occurred while reading message: %v", err)
			ws.Close()
			delete(clients, ws)
			break
		}
		if datatype == websocket.CloseMessage {
			log.Println("CloseMessage received")
			delete(clients, ws)
			break
		}
		var datas = &BufferReader{
			Index: 0,
			Bytes: data,
		}
		datas.Buff.Write(data)
		if !datas.Eof() {
			cmd, err := datas.ReadCString()
			if err != nil {
				panic(err)
			}
			// fmt.Printf("CMD : [%s]\n", cmd)
			switch cmd {
			case "hello":
				h.handleRequest(cmd)
			case "version":
				version, err := datas.ReadUInt32LE()
				if err != nil {
					panic(err)
				}
				fmt.Printf("VERSION : %d", version)
				h.handleRequest(cmd)
			case "dataStop":
				h.handleRequest(cmd)
			case "dataStart":
				h.handleRequest(cmd)
			case "levelInit":
				mapname, err := datas.ReadCString()
				if err != nil {
					panic(err)
				}
				fmt.Printf("map : %s", mapname) //
				h.handleRequest(cmd)
			case "levelShutdown":
				h.handleRequest(cmd)
			case "cam":
				time, err := datas.ReadFloatLE()
				xPosition, err := datas.ReadFloatLE()
				yPosition, err := datas.ReadFloatLE()
				zPosition, err := datas.ReadFloatLE()
				xRotation, err := datas.ReadFloatLE()
				yRotation, err := datas.ReadFloatLE()
				zRotation, err := datas.ReadFloatLE()
				fov, err := datas.ReadFloatLE()
				if err != nil {
					panic(err)
				}
				// fmt.Printf("time = %f\n", time)
				// fmt.Printf("xPosition = %f\n", xPosition)
				// fmt.Printf("yPosition = %f\n", yPosition)
				// fmt.Printf("zPosition = %f\n", zPosition)
				// fmt.Printf("xRotation = %f\n", xRotation)
				// fmt.Printf("yRotation = %f\n", yRotation)
				// fmt.Printf("zRotation = %f\n", zRotation)
				// fmt.Printf("fov = %f\n", fov)
				cam := &CamData{
					XPos: xPosition,
					YPos: yPosition,
					ZPos: zPosition,
					XRot: xRotation,
					YRot: yRotation,
					Zrot: zRotation,
					Time: time,
					Fov:  fov,
				}
				h.handleCamRequest(cam)
				//
			case "gameEvent":
				//TODO. JSON
				h.handleRequest(cmd)
			default:
				fmt.Println("Unknown message")
				h.handleRequest(cmd)
			}
		}
	}
}
