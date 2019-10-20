package mirvpgl

import (
	"math"
	"strings"
	"unsafe"
	//"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	//"log"
	"github.com/gorilla/websocket"
	"math/big"
	//"net/http"
	// "strconv"
)

var (
	nullstr byte = []byte("\x00")[0]
)

func FindDelim(buffer []byte, idx int) int { // byte?? 区切り文字の位置を返すfunc
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
	val, err := b.Buff.ReadBytes(byte(b.Index))
	if err != nil {
		return 0 // todo
	}
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
	fmt.Printf("\nb.index : %d\nb.bytes len : %d\n", b.Index, b.Buff.Len())
	if b.Index >= b.Buff.Len() {
		fmt.Println("EOF")
		return true
	}
	return false
}

func SendRCON(ws *websocket.Conn, cmd string) {
	command := []byte("exec")
	command = append(command, nullstr)
	command = append(command, []byte(cmd)...)
	command = append(command, nullstr)
	ws.WriteMessage(2, []uint8(command))
}
