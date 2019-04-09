package env

import (
	"fmt"

	"github.com/tanema/mal/wotlisp/src/types"
)

type Env struct {
	data  map[string]types.Base
	outer *Env
}

func New(outer *Env) *Env {
	env := &Env{data: map[string]types.Base{}, outer: outer}
	return env
}

func (e *Env) Find(key types.Symbol) *Env {
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
	return env.data[string(key)], nil
}
