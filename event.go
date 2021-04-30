package mirvpgl

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"strconv"
)

const (
	_ int32 = iota
	KEYTYPE_STRING
	KEYTYPE_FLOAT32
	KEYTYPE_INT32
	KEYTYPE_INT16
	KEYTYPE_INT8
	KEYTYPE_BOOLEAN
	KEYTYPE_BIGUINT64
	KEYTYPE_UNKNOWN
)

type gameEventUnserializer struct {
	Enrichments Enrichments
	KnownEvents map[int32]*GameEventDescription // id->event desc
}

func newGameEventUnserializer(e Enrichments) *gameEventUnserializer {
	return &gameEventUnserializer{
		Enrichments: e,
		KnownEvents: make(map[int32]*GameEventDescription, 0),
	}
}

func (g *gameEventUnserializer) Unserialize(r io.Reader) (*GameEventData, error) {
	var ev *GameEventDescription
	var eventID int32
	if err := binary.Read(r, binary.LittleEndian, &eventID); err != nil {
		return nil, err
	}
	if eventID == 0 {
		gameEvent, err := newGameEventDescription(r)
		if err != nil {
			return nil, err
		}
		g.KnownEvents[gameEvent.EventID] = gameEvent

		if _, ok := g.Enrichments[gameEvent.EventName]; ok {
			gameEvent.enrichments = g.Enrichments[gameEvent.EventName]
		}
		ev = gameEvent
	} else {
		e, ok := g.KnownEvents[eventID]
		if !ok {
			ev = &GameEventDescription{}
		} else {
			ev = e
		}
	}
	fmt.Println("ev:", ev)
	return ev.Unserialize(r)
}

// Cordinates include float32 X/Y/Z Pos cordinates.
type Cordinates struct {
	X float32
	Y float32
	Z float32
}

// GameEventDescription include Event ID,Name, Keys etc.
type GameEventDescription struct {
	EventID     int32
	EventName   string
	Keys        []EventKey // KeyName->Key type
	enrichments map[string]Enrichment
	// enrichment // see https://wiki.alliedmods.net/Counter-Strike:_Global_Offensive_Events
}

func newGameEventDescription(r io.Reader) (*GameEventDescription, error) {
	d := &GameEventDescription{
		EventID:   0,
		EventName: "",
		Keys:      make([]EventKey, 0),
	}
	buf := bufio.NewReader(r)
	if err := binary.Read(r, binary.LittleEndian, &d.EventID); err != nil {
		return nil, fmt.Errorf("Failed to parse Event ID : %v", err)
	}

	eventName, err := buf.ReadString(nullstr)
	if err != nil {
		return nil, err
	}
	d.EventName = eventName
	for {
		if ok, err := buf.ReadByte(); err != nil {
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return nil, fmt.Errorf("Failed to read ok value:%v", err)
				}
			}
			if ok == 0 {
				break
			}
			keyName, err := buf.ReadString(nullstr)
			if err != nil {
				return nil, err
			}
			var keyType int32
			if err := binary.Read(r, binary.LittleEndian, &keyType); err != nil {
				return nil, err
			}
			d.Keys = append(d.Keys, EventKey{
				Name: keyName,
				Type: keyType,
			})
		}
	}

	return d, nil
}

// GameEventData Game event keys and time
type GameEventData struct {
	Name       string
	ClientTime float32
	Keys       map[string]string // Even value is float32 or int etc. convert to string
}

// Unserialize parse EventDescription
func (e *GameEventDescription) Unserialize(r io.Reader) (*GameEventData, error) {
	d := &GameEventData{}

	buf := bufio.NewReader(r)
	if err := binary.Read(r, binary.LittleEndian, &d.ClientTime); err != nil {
		return nil, fmt.Errorf("Failed to read client time:%v", err)
	}

	for _, v := range e.Keys {
		keyname := v.Name
		var keyvalue string
		switch v.Type {
		case KEYTYPE_STRING:
			val, err := buf.ReadString(nullstr)
			if err != nil {
				return nil, fmt.Errorf("Failed to read CString value:%v", err)
			}
			keyvalue = val
		case KEYTYPE_FLOAT32:
			var f float32
			if err := binary.Read(r, binary.LittleEndian, &f); err != nil {
				return nil, fmt.Errorf("Failed to read float32 value:%v", err)
			}
			keyvalue = strconv.FormatFloat(float64(f), 'f', -1, 64)
		case KEYTYPE_INT32:
			var f int32
			if err := binary.Read(r, binary.LittleEndian, &f); err != nil {
				return nil, fmt.Errorf("Failed to read int32 value:%v", err)
			}
			keyvalue = fmt.Sprint(f)
		case KEYTYPE_INT16:
			var f int16
			if err := binary.Read(r, binary.LittleEndian, &f); err != nil {
				return nil, fmt.Errorf("Failed to read int16 value:%v", err)
			}
			keyvalue = fmt.Sprint(f)
		case KEYTYPE_INT8:
			var f int8
			if err := binary.Read(r, binary.LittleEndian, &f); err != nil {
				return nil, fmt.Errorf("Failed to read int8 value:%v", err)
			}
			keyvalue = fmt.Sprint(f)
		case KEYTYPE_BOOLEAN:
			var f bool
			if err := binary.Read(r, binary.LittleEndian, &f); err != nil {
				return nil, fmt.Errorf("Failed to read boolean value:%v", err)
			}
			keyvalue = fmt.Sprint(f)
		case KEYTYPE_BIGUINT64:
			var f1 uint32
			var f2 uint32
			if err := binary.Read(r, binary.LittleEndian, &f1); err != nil {
				return nil, fmt.Errorf("Failed to read bigint64 value:%v", err)
			}
			if err := binary.Read(r, binary.LittleEndian, &f2); err != nil {
				return nil, fmt.Errorf("Failed to read bigint64 value:%v", err)
			}
			var lo *big.Int
			var hi *big.Int
			lo = lo.SetUint64(uint64(f1))
			hi = hi.SetUint64(uint64(f2))
			var f *big.Int
			keyvalue = f.Or(lo, hi.Lsh(hi, 32)).String()
		default:
			return nil, fmt.Errorf("Unknown Event key")
		}
		d.Keys[keyname] = keyvalue
		// Check enrichments keyName check...
	}
	return d, nil
}

// EventKey key-value struct with dynamic typing
type EventKey struct {
	Name string
	Type int32
}
