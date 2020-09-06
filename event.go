package mirvpgl

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
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

// Cordinates include float32 X/Y/Z Pos cordinates.
type Cordinates [3]float32

// EventDescription include Event ID,Name, Keys etc.
type EventDescription struct {
	EventID    int32
	EventName  string
	Keys       []EventKeys
	ClientTime float32
}

// EventKeys key-value struct with dynamic typing
type EventKeys struct {
	Name  string
	Type  int32
	Value interface{} // TODO...
}

// ParseEvent parse EventDescription
func ParseEvent(r io.Reader, desc *EventDescription) error {
	if err := binary.Read(r, binary.LittleEndian, desc.EventID); err != nil {
		return fmt.Errorf("Failed to parse Event ID : %v", err)
	}
	if err := binary.Read(r, binary.LittleEndian, desc.EventName); err != nil {
		return fmt.Errorf("Failed to parse Event name : %v", err)
	}
	if desc.Keys == nil {
		desc.Keys = make([]EventKeys, 0)
	}
	var ok bool
	for ok {
		if err := binary.Read(r, binary.LittleEndian, &ok); err != nil {
			return err
		}
		var keyName string
		if err := binary.Read(r, binary.LittleEndian, &keyName); err != nil {
			return fmt.Errorf("Failed to parse key name : %v", err)
		}
		var keyType int32
		if err := binary.Read(r, binary.LittleEndian, &keyType); err != nil {
			return fmt.Errorf("Failed to parse key name : %v", err)
		}
		desc.Keys = append(desc.Keys, EventKeys{
			Name:  keyName,
			Type:  keyType,
			Value: nil,
		})
	}
	if err := binary.Read(r, binary.LittleEndian, desc.ClientTime); err != nil {
		return fmt.Errorf("Failed to parse Client time : %v", err)
	}
	for k, v := range desc.Keys {
		switch v.Type {
		case KEYTYPE_STRING:
			var v string
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				return fmt.Errorf("Failed to parse Key value : %v", err)
			}
			desc.Keys[k].Value = v
		case KEYTYPE_FLOAT32:
			var v float32
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				return fmt.Errorf("Failed to parse Key value : %v", err)
			}
			desc.Keys[k].Value = v
		case KEYTYPE_INT32:
			var v int32
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				return fmt.Errorf("Failed to parse Key value : %v", err)
			}
			desc.Keys[k].Value = v
		case KEYTYPE_INT16:
			var v int16
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				return fmt.Errorf("Failed to parse Key value : %v", err)
			}
			desc.Keys[k].Value = v
		case KEYTYPE_INT8:
			var v int8
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				return fmt.Errorf("Failed to parse Key value : %v", err)
			}
			desc.Keys[k].Value = v
		case KEYTYPE_BOOLEAN:
			var v bool
			if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
				return fmt.Errorf("Failed to parse Key value : %v", err)
			}
			desc.Keys[k].Value = v
		case KEYTYPE_BIGUINT64:
			ar := [8]byte{}
			if err := binary.Read(r, binary.LittleEndian, &ar); err != nil {
				return fmt.Errorf("Failed to parse Key value : %v", err)
			}
			v := big.Int{}
			sl := make([]byte, 0, 8)
			copy(sl, ar[:])
			v.SetBytes(sl)
			desc.Keys[k].Value = v
		default:
			return fmt.Errorf("Failed to parse Key type")
		}
	}
	return nil
}
