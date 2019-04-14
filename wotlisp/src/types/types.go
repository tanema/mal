package types

import (
	"errors"
	"fmt"
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

func NewHashmap(values []Base, excludeKeys ...Base) (*Hashmap, error) {
	if len(values)%2 == 1 {
		return nil, errors.New("Odd number of arguments to NewHashMap")
	}
	m := map[Base]Base{}
	for i := 0; i < len(values); i += 2 {
		key := values[i]
		found := false
		for _, exclude := range excludeKeys {
			if key == exclude {
				found = true
				break
			}
		}
		if !found {
			m[key] = values[i+1]
		}
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

func (hm *Hashmap) Keys() []Base {
	keys := make([]Base, 0, len(hm.Forms))
	for k := range hm.Forms {
		keys = append(keys, k)
	}
	return keys
}

func (hm *Hashmap) Vals() []Base {
	vals := make([]Base, 0, len(hm.Forms))
	for _, v := range hm.Forms {
		vals = append(vals, v)
	}
	return vals
}

type ExtFunc struct {
	AST     Base
	Params  []Base
	Env     Env
	IsMacro bool
	eval    func(Env, Base) (Base, error)
}

func NewFunc(env Env, eval func(Env, Base) (Base, error), args ...Base) (*ExtFunc, error) {
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
	return fn.eval(newEnv, fn.AST)
}

func CallFunc(e Env, baseFn Base, arguments []Base) (Base, error) {
	switch fn := baseFn.(type) {
	case Func:
		return fn(e, arguments)
	case *ExtFunc:
		return fn.Apply(arguments)
	default:
		return nil, fmt.Errorf("attempt to call non-function %v", baseFn)
	}
}

type UserError struct {
	Val Base
}

func (err UserError) Error() string {
	return "User Error"
}
