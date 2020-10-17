package mirvpgl

import (
	"encoding/binary"
	"fmt"
	"io"
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
	keyvalue   []EventKeys
	Keys       map[string]interface{}
	ClientTime float32
}

// EventKeys key-value struct with dynamic typing
type EventKeys struct {
	Name string
	Type int32
}

// ParseEvent parse EventDescription
func (e *EventDescription) ParseEvent(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &e.EventID); err != nil {
		return fmt.Errorf("Failed to parse Event ID : %v", err)
	}
	if e.EventID == 0 {
		if err := binary.Read(r, binary.LittleEndian, &e.EventID); err != nil {
			return fmt.Errorf("Failed to parse Event ID : %v", err)
		}
		e.EventName, err = buf.ReadString(nullstr)
		if err != nil {
			return err
		}
		for {
			var ok bool
			if err := binary.Read(r, binary.LittleEndian, &ok); err != nil {
				return err
			}
			if !ok {
				break
			}
			keyName, err := buf.ReadString(nullstr)
			if err != nil {
				return err
			}
			var keyType int32
			if err := binary.Read(r, binary.LittleEndian, &keyType); err != nil {
				return err
			}
			e.keyvalue = append(e.keyvalue, EventKeys{
				Name: keyName,
				Type: keyType,
			})
		}
	}
	if err := binary.Read(r, binary.LittleEndian, &e.ClientTime); err != nil {
		return err
	}
	for k, v := range e.keyvalue {
		key := v
		keyName := v.Name
		var keyValue interface{}

		switch key.Type {
		case KEYTYPE_STRING:
			keyValue, err = buf.ReadString(nullstr)
			if err != nil {
				return err
			}
		case KEYTYPE_FLOAT32:
			var f float32
			if err := binary.Read(&f); err != nil {
				return err
			}
			keyValue = f
		case KEYTYPE_INT32:
			var f int32
			if err := binary.Read(&f); err != nil {
				return err
			}
			keyValue = f
		case KEYTYPE_INT16:
			var f int16
			if err := binary.Read(&f); err != nil {
				return err
			}
			keyValue = f
		case KEYTYPE_INT16:
			var f int16
			if err := binary.Read(&f); err != nil {
				return err
			}
			keyValue = f
		case KEYTYPE_INT8:
			var f int8
			if err := binary.Read(&f); err != nil {
				return err
			}
			keyValue = f
		case KEYTYPE_BOOLEAN:
			var f bool
			if err := binary.Read(&f); err != nil {
				return err
			}
			keyValue = f
		case KEYTYPE_BIGUINT64:
			// ??
		default:
			return fmt.Errorf("Unknown Event key")
		}
		// Check enrichments keyName check...
	}
	return nil
}
