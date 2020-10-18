package mirvpgl

import (
	"bufio"
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
	keyvalue   []EventKeys
	Keys       map[string]interface{}
	ClientTime float32
	// enrichment // see https://wiki.alliedmods.net/Counter-Strike:_Global_Offensive_Events
}

// Unserialize parse EventDescription
func (e *EventDescription) Unserialize(r io.Reader) error {
	buf := bufio.NewReader(r)
	if err := binary.Read(r, binary.LittleEndian, &e.EventID); err != nil {
		return fmt.Errorf("Failed to parse Event ID : %v", err)
	}
	if e.EventID == 0 {
		if err := binary.Read(r, binary.LittleEndian, &e.EventID); err != nil {
			return fmt.Errorf("Failed to parse Event ID : %v", err)
		}
		var err error
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
	e.Keys = make(map[string]interface{})

	for _, v := range e.keyvalue {
		key := v
		keyName := v.Name
		var keyValue interface{}

		switch key.Type {
		case KEYTYPE_STRING:
			var err error
			keyValue, err = buf.ReadString(nullstr)
			if err != nil {
				return err
			}
		case KEYTYPE_FLOAT32:
			var f float32
			if err := binary.Read(r, binary.LittleEndian, &f); err != nil {
				return err
			}
			keyValue = f
		case KEYTYPE_INT32:
			var f int32
			if err := binary.Read(r, binary.LittleEndian, &f); err != nil {
				return err
			}
			keyValue = f
		case KEYTYPE_INT16:
			var f int16
			if err := binary.Read(r, binary.LittleEndian, &f); err != nil {
				return err
			}
			keyValue = f
		case KEYTYPE_INT8:
			var f int8
			if err := binary.Read(r, binary.LittleEndian, &f); err != nil {
				return err
			}
			keyValue = f
		case KEYTYPE_BOOLEAN:
			var f bool
			if err := binary.Read(r, binary.LittleEndian, &f); err != nil {
				return err
			}
			keyValue = f
		case KEYTYPE_BIGUINT64:
			var f1 uint32
			var f2 uint32
			if err := binary.Read(r, binary.LittleEndian, &f1); err != nil {
				return err
			}
			if err := binary.Read(r, binary.LittleEndian, &f2); err != nil {
				return err
			}
			var lo *big.Int
			var hi *big.Int
			lo = lo.SetUint64(uint64(f1))
			hi = hi.SetUint64(uint64(f2))
			var f *big.Int
			keyValue = f.Or(lo, hi.Lsh(hi, 32)).String()
		default:
			return fmt.Errorf("Unknown Event key")
		}
		e.Keys[keyName] = keyValue
		// Check enrichments keyName check...
	}
	return nil
}

// EventKeys key-value struct with dynamic typing
type EventKeys struct {
	Name string
	Type int32
}

// UserIDEnrichment contains User informations with XUID/Eyeorigins(Cordinates)/EyeAngles(Cordinates)
type UserIDEnrichment struct {
	XUID      string
	EyeOrigin Cordinates
	EyeAngles Cordinates
	// keyValue??
}

// Unserialize UserID Enrichment
func (u *UserIDEnrichment) Unserialize(r io.Reader) error {
	var f1 uint32
	var f2 uint32
	if err := binary.Read(r, binary.LittleEndian, &f1); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &f2); err != nil {
		return err
	}
	var u1 *big.Int
	var u2 *big.Int
	u1 = u1.SetUint64(uint64(f1))
	u2 = u2.SetUint64(uint64(f2))
	var f *big.Int
	u.XUID = f.Add(u1, u2).String()

	if err := binary.Read(r, binary.LittleEndian, &u.EyeOrigin); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &u.EyeAngles); err != nil {
		return err
	}
	return nil
}

// EntityNumEnrichment containns Entity's Origin/Angles.
type EntityNumEnrichment struct {
	Origin Cordinates
	Angles Cordinates
	// KeyValue??
}

// Unserialize EntityNum Enrichment
func (e *EntityNumEnrichment) Unserialize(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &e.Origin); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &e.Angles); err != nil {
		return err
	}
	return nil
}
