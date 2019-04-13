package types

import (
	"errors"
)

type (
	Base    interface{}
	Atom    struct{ Val Base }
	Symbol  string
	Keyword string
	Func    func(Env, []Base) (Base, error)
)

type Env interface {
	Child([]Base, []Base) (Env, error)
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
	AST     Base
	Params  []Base
	Env     Env
	IsMacro bool
	eval    func(Base, Env) (Base, error)
}

func NewFunc(env Env, eval func(Base, Env) (Base, error), args ...Base) (*ExtFunc, error) {
	if len(args) < 2 {
		return nil, errors.New("improperly formatted fn* statement")
	}

	var params []Base
	switch tparams := args[0].(type) {
	case Collection:
		params = tparams.Data()
	default:
		return nil, errors.New("invalid fn* param declaration")
	}

	return &ExtFunc{
		AST:    args[1],
		Params: params,
		Env:    env,
		eval:   eval,
	}, nil
}

func (fn *ExtFunc) Apply(arguments []Base) (Base, error) {
	newEnv, err := fn.Env.Child(fn.Params, arguments)
	if err != nil {
		return nil, err
	}
	return fn.eval(fn.AST, newEnv)
}
