package types

import (
	"errors"
)

type (
	Base    interface{}
	Symbol  string
	Keyword string
)

type EnvType interface {
	Find(key Symbol) EnvType
	Set(key Symbol, value Base) Base
	Get(key Symbol) (Base, error)
}

type List struct {
	Forms []Base
}

func NewList(forms ...Base) *List {
	return &List{Forms: forms}
}

type Vector struct {
	Forms []Base
}

type Hashmap struct {
	Forms map[Base]Base
}

func NewHashmap(values []Base) (*Hashmap, error) {
	if len(values)%2 == 1 {
		return nil, errors.New("Odd number of arguments to NewHashMap")
	}
	m := map[Base]Base{}
	for i := 0; i < len(values); i += 2 {
		m[values[i]] = values[i+1]
	}
	return &Hashmap{Forms: m}, nil
}

func (hm *Hashmap) ToList() []Base {
	values := []Base{}
	for key, val := range hm.Forms {
		values = append(values, key)
		values = append(values, val)
	}
	return values
}
