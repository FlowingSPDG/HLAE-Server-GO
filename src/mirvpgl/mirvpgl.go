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
	"math/big"
	//"net/http"
	// "strconv"
	//"github.com/gorilla/websocket"
)

var (
	nullstr byte = []byte("\x00")[0]
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
