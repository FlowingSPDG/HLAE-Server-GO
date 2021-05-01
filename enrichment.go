package mirvpgl

import (
	"encoding/binary"
	"io"
	"math/big"
)

type Enrichments map[string]map[string]Enrichment

type Enrichment interface {
	Unserialize(r io.Reader) error
	GetMap() map[string]interface{}
	SetEnrichment([]string)
	GetEnrichment() []string
}

// UserIDEnrichment contains User informations with XUID/Eyeorigins(Cordinates)/EyeAngles(Cordinates)
type UserIDEnrichment struct {
	enrichments []string
	XUID        *big.Int
	EyeOrigin   Cordinates
	EyeAngles   Cordinates
	KeyValue    string
}

func newUserIDEnrichment() *UserIDEnrichment {
	return &UserIDEnrichment{
		enrichments: []string{
			"useridWithSteamId",
			"useridWithEyePosition",
			"useridWithEyeAngles",
		},
		XUID:      &big.Int{},
		EyeOrigin: Cordinates{},
		EyeAngles: Cordinates{},
		KeyValue:  "",
	}
}

// Unserialize Unserialize into u
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
	u.XUID = f.Add(u1, u2)

	if err := binary.Read(r, binary.LittleEndian, &u.EyeOrigin); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &u.EyeAngles); err != nil {
		return err
	}
	return nil
}

// GetMap UserID Enrichment to map
func (u *UserIDEnrichment) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"value":     u.KeyValue,
		"xuid":      u.XUID.String(),
		"eyeOrigin": u.EyeOrigin,
		"eyeAngles": u.EyeAngles,
	}
}

func (u *UserIDEnrichment) SetEnrichment(en []string) {
	u.enrichments = en
}

func (u *UserIDEnrichment) GetEnrichment() []string {
	return u.enrichments
}

// EntityNumEnrichment containns Entity's Origin/Angles.
type EntityNumEnrichment struct {
	enrichments []string
	Origin      Cordinates
	Angles      Cordinates
	KeyValue    string
}

func newEntityNumEnrichment() *EntityNumEnrichment {
	return &EntityNumEnrichment{
		enrichments: []string{
			"entnumWithOrigin",
			"entnumWithAngles",
		},
		Origin:   Cordinates{},
		Angles:   Cordinates{},
		KeyValue: "",
	}
}

// GetMap EntityNum Enrichment
func (e *EntityNumEnrichment) Unserialize(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &e.Origin); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &e.Angles); err != nil {
		return err
	}
	return nil
}

// GetMap EntityNum Enrichment
func (e *EntityNumEnrichment) GetMap() map[string]interface{} {
	return map[string]interface{}{
		"value":  e.KeyValue,
		"origin": e.Origin,
		"angles": e.Angles,
	}
}

func (e *EntityNumEnrichment) SetEnrichment(en []string) {
	e.enrichments = en
}
func (e *EntityNumEnrichment) GetEnrichment() []string {
	return e.enrichments
}
