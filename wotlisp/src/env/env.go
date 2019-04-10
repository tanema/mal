package env

import (
	"fmt"

	"github.com/tanema/mal/wotlisp/src/types"
)

type Env struct {
	data  map[string]types.Base
	outer types.Env
}

func New(outer types.Env, binds, exprs []types.Base) (*Env, error) {
	env := &Env{data: map[string]types.Base{}, outer: outer}
	for i, bind := range binds {
		key, ok := bind.(types.Symbol)
		if !ok {
			return nil, fmt.Errorf("non-symbol bind value")
		}
		if key == "&" {
			env.Set(binds[i+1].(types.Symbol), types.NewList(exprs[i:]...))
			break
		}
		env.Set(key, exprs[i])
	}
	return env, nil
}

func (e *Env) Find(key types.Symbol) types.Env {
	if _, ok := e.data[string(key)]; ok {
		return e
	} else if e.outer != nil {
		return e.outer.Find(key)
	}
	return nil
}

func (e *Env) Set(key types.Symbol, value types.Base) {
	e.data[string(key)] = value
}

func (e *Env) Get(key types.Symbol) (types.Base, error) {
	env := e.Find(key)
	if env == nil {
		return nil, fmt.Errorf("variable %v not found", key)
	}
	return env.(*Env).data[string(key)], nil
}
