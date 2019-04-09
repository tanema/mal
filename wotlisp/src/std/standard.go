package std

import (
	"errors"

	"github.com/tanema/mal/wotlisp/src/types"
)

var StdEnv = map[string]types.Base{
	"+": func(a []types.Base) (types.Base, error) {
		if e := assertArgNum(a, 2); e != nil {
			return nil, e
		}
		return a[0].(int) + a[1].(int), nil
	},
	"-": func(a []types.Base) (types.Base, error) {
		if e := assertArgNum(a, 2); e != nil {
			return nil, e
		}
		return a[0].(int) - a[1].(int), nil
	},
	"*": func(a []types.Base) (types.Base, error) {
		if e := assertArgNum(a, 2); e != nil {
			return nil, e
		}
		return a[0].(int) * a[1].(int), nil
	},
	"/": func(a []types.Base) (types.Base, error) {
		if e := assertArgNum(a, 2); e != nil {
			return nil, e
		}
		return a[0].(int) / a[1].(int), nil
	},
}

func assertArgNum(a []types.Base, expectedLen int) error {
	if len(a) != expectedLen {
		return errors.New("wrong number of arguments")
	}
	return nil
}
