package types

import (
	"errors"
)

type (
	Base    interface{}
	Symbol  string
	Keyword string
	Func    func(Env, []Base) (Base, error)
)

type Env interface {
	Find(Symbol) Env
	Set(Symbol, Base)
	Get(Symbol) (Base, error)
}

type Collection interface {
	Data() []Base
}

type List struct {
	Forms []Base
}

func NewList(forms ...Base) *List {
	return &List{Forms: forms}
}

func (l *List) Data() []Base { return l.Forms }

type Vector struct {
	Forms []Base
}

func (l *Vector) Data() []Base { return l.Forms }

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

type ExtFunc struct {
	AST    Base
	Params []Base
	Env    Env
	Fn     Func
}
