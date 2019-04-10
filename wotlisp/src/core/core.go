package core

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/tanema/mal/wotlisp/src/printer"
	"github.com/tanema/mal/wotlisp/src/types"
)

var Namespace = map[types.Symbol]types.Func{
	"+":       add,
	"-":       sub,
	"*":       mul,
	"/":       div,
	"=":       equal,
	"<":       lessThan,
	"<=":      lessThanEqual,
	">":       greaterThan,
	">=":      greaterThanEqual,
	"prn":     prn,
	"println": prnln,
	"pr-str":  prnstr,
	"str":     str,
	"list":    list,
	"list?":   islist,
	"empty?":  isempty,
	"count":   count,
}

func prn(e types.Env, a []types.Base) (types.Base, error) {
	fmt.Println(printer.List(a, true, "", "", " "))
	return nil, nil
}

func prnln(e types.Env, a []types.Base) (types.Base, error) {
	fmt.Println(printer.List(a, false, "", "", " "))
	return nil, nil
}

func prnstr(e types.Env, a []types.Base) (types.Base, error) {
	return printer.List(a, true, "", "", " "), nil
}

func str(e types.Env, a []types.Base) (types.Base, error) {
	return printer.List(a, false, "", "", ""), nil
}

func list(e types.Env, a []types.Base) (types.Base, error) {
	return types.NewList(a...), nil
}

func islist(e types.Env, a []types.Base) (types.Base, error) {
	if len(a) == 0 {
		return false, errors.New("not enough arguments to list?")
	}
	_, lst := a[0].(*types.List)
	return lst, nil
}

func isempty(e types.Env, a []types.Base) (types.Base, error) {
	if len(a) == 0 {
		return false, errors.New("not enough arguments to empty?")
	}
	switch data := a[0].(type) {
	case types.Collection:
		return len(data.Data()) == 0, nil
	case *types.Hashmap:
		return len(data.Forms) == 0, nil
	case nil:
		return true, nil
	default:
		return false, errors.New("invalid data type passed to empty?")
	}
}

func count(e types.Env, a []types.Base) (types.Base, error) {
	if len(a) == 0 {
		return false, errors.New("nothing to count")
	}
	switch data := a[0].(type) {
	case types.Collection:
		return float64(len(data.Data())), nil
	case *types.Hashmap:
		return float64(len(data.Forms)), nil
	case string:
		return float64(len(data)), nil
	case nil:
		return float64(0), nil
	default:
		return false, errors.New("invalid data type passed to count")
	}
}

func equal(e types.Env, a []types.Base) (types.Base, error) {
	if len(a) < 2 {
		return false, errors.New("not enough arguments to equal")
	}
	return checkEquality(a[0], a[1])
}

func checkEquality(val1, val2 types.Base) (bool, error) {
	switch data := val1.(type) {
	case types.Collection:
		other, ok := val2.(types.Collection)
		if !ok {
			return false, nil
		}
		return equalLists(data.Data(), other.Data())
	default:
	}

	if reflect.TypeOf(val1) != reflect.TypeOf(val2) {
		return false, nil
	}

	switch data := val1.(type) {
	case *types.Hashmap:
		other := val2.(*types.Hashmap)
		return equalMaps(data.Forms, other.Forms)
	case types.Keyword:
		other := val2.(types.Keyword)
		return data == other, nil
	case string:
		other := val2.(string)
		return data == other, nil
	case float64:
		other := val2.(float64)
		return data == other, nil
	case bool:
		other := val2.(bool)
		return data == other, nil
	case nil:
		return true, nil
	default:
		return false, errors.New("invalid data type passed to equal")
	}
}

func equalLists(lst1, lst2 []types.Base) (bool, error) {
	if len(lst1) != len(lst2) {
		return false, nil
	}

	for i, elm := range lst1 {
		if equal, err := checkEquality(elm, lst2[i]); err != nil {
			return false, err
		} else if !equal {
			return false, nil
		}
	}

	return true, nil
}

func equalMaps(m1, m2 map[types.Base]types.Base) (bool, error) {
	if len(m1) != len(m2) {
		return false, nil
	}

	for key, val := range m1 {
		other, found := m2[key]
		if !found {
			return false, nil
		}

		if equal, err := checkEquality(val, other); err != nil {
			return false, err
		} else if !equal {
			return false, nil
		}
	}

	return true, nil
}

func prepareCompare(args []types.Base) (float64, float64, error) {
	if len(args) < 2 {
		return 0, 0, errors.New("not enough arguments to equal")
	}
	val1, ok := args[0].(float64)
	if !ok {
		return 0, 0, errors.New("cannot compare non-number values")
	}
	val2, ok := args[1].(float64)
	if !ok {
		return 0, 0, errors.New("cannot compare non-number values")
	}
	return val1, val2, nil
}

func lessThan(e types.Env, a []types.Base) (types.Base, error) {
	v1, v2, err := prepareCompare(a)
	if err != nil {
		return nil, err
	}
	return v1 < v2, nil
}

func lessThanEqual(e types.Env, a []types.Base) (types.Base, error) {
	v1, v2, err := prepareCompare(a)
	if err != nil {
		return nil, err
	}
	return v1 <= v2, nil
}

func greaterThan(e types.Env, a []types.Base) (types.Base, error) {
	v1, v2, err := prepareCompare(a)
	if err != nil {
		return nil, err
	}
	return v1 > v2, nil
}

func greaterThanEqual(e types.Env, a []types.Base) (types.Base, error) {
	v1, v2, err := prepareCompare(a)
	if err != nil {
		return nil, err
	}
	return v1 >= v2, nil
}

func add(e types.Env, a []types.Base) (types.Base, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(float64) + a[1].(float64), nil
}

func sub(e types.Env, a []types.Base) (types.Base, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(float64) - a[1].(float64), nil
}

func mul(e types.Env, a []types.Base) (types.Base, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(float64) * a[1].(float64), nil
}

func div(e types.Env, a []types.Base) (types.Base, error) {
	if e := assertArgNum(a, 2); e != nil {
		return nil, e
	}
	return a[0].(float64) / a[1].(float64), nil
}

func assertArgNum(a []types.Base, expectedLen int) error {
	if len(a) != expectedLen {
		return errors.New("wrong number of arguments")
	}
	return nil
}
