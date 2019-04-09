package std

import (
	"errors"

	"github.com/tanema/mal/wotlisp/src/env"
	"github.com/tanema/mal/wotlisp/src/types"
)

func Define(e *env.Env) {
	e.Set("+", add)
	e.Set("-", sub)
	e.Set("*", mul)
	e.Set("/", div)
}

func add(e *env.Env, a []types.Base) (types.Base, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(int) + a[1].(int), nil
}

func sub(e *env.Env, a []types.Base) (types.Base, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(int) - a[1].(int), nil
}

func mul(e *env.Env, a []types.Base) (types.Base, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(int) * a[1].(int), nil
}

func div(e *env.Env, a []types.Base) (types.Base, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(int) / a[1].(int), nil
}

func assertArgNum(a []types.Base, expectedLen int) error {
	if len(a) != expectedLen {
		return errors.New("wrong number of arguments")
	}
	return nil
}
