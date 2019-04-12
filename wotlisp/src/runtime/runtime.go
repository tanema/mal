package runtime

import (
	"fmt"

	"github.com/tanema/mal/wotlisp/src/env"
	"github.com/tanema/mal/wotlisp/src/types"
)

func Eval(object types.Base, e types.Env) (types.Base, error) {
	var err error
	for {
		switch tobject := object.(type) {
		case *types.List:
			if len(tobject.Forms) == 0 {
				return tobject, nil
			}
			sym, _ := tobject.Forms[0].(types.Symbol)
			switch sym {
			case "quote":
				if len(tobject.Forms) < 2 {
					return nil, nil
				}
				return tobject.Forms[1], nil
			case "quasiquote":
				object = evalQuasiQuote(e, tobject.Forms[1])
			case "do":
				object, err = evalDo(e, tobject.Forms[1:]...)
			case "if":
				object, err = evalIf(e, tobject.Forms[1:]...)
			case "fn*":
				return evalFn(e, tobject.Forms[1:]...)
			case "def!":
				return evalDef(e, tobject.Forms[1:]...)
			case "let*":
				object, e, err = evalLet(e, tobject.Forms[1:]...)
			default:
				lst, err := evalAST(tobject, e)
				if err != nil {
					return nil, err
				}
				list := lst.(*types.List)
				switch fn := list.Forms[0].(type) {
				case types.Func:
					return fn(e, list.Forms[1:])
				case *types.ExtFunc:
					newEnv, err := env.New(fn.Env, fn.Params, list.Forms[1:])
					if err != nil {
						return nil, err
					}
					object, e = fn.AST, newEnv
				default:
					return nil, fmt.Errorf("attempt to call non-function %v", list.Forms[0])
				}
			}
		default:
			return evalAST(tobject, e)
		}

		if err != nil {
			return nil, err
		}
	}
}

func evalAST(ast types.Base, env types.Env) (types.Base, error) {
	switch tobject := ast.(type) {
	case types.Symbol:
		symVal, err := env.Get(tobject)
		if err != nil {
			return nil, err
		}
		return symVal, nil
	case *types.List:
		lst, err := evalListForms(tobject.Data(), env)
		return &types.List{Forms: lst}, err
	case *types.Vector:
		lst, err := evalListForms(tobject.Forms, env)
		return &types.Vector{Forms: lst}, err
	case *types.Hashmap:
		lst, err := evalListForms(tobject.ToList(), env)
		if err != nil {
			return nil, err
		}
		return types.NewHashmap(lst)
	default:
		return ast, nil
	}
}

func evalDef(e types.Env, args ...types.Base) (types.Base, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments")
	}

	name, ok := args[0].(types.Symbol)
	if !ok {
		return nil, fmt.Errorf("non-symbol bind value")
	}
	value, err := Eval(args[1], e)
	if err == nil {
		e.Set(name, value)
	}
	return value, err
}

func evalLet(e types.Env, args ...types.Base) (types.Base, types.Env, error) {
	if len(args) < 2 {
		return nil, nil, fmt.Errorf("not enough arguments for  let* call")
	}
	newEnv, err := env.New(e, nil, nil)
	if err != nil {
		return nil, nil, err
	}

	var definitions []types.Base
	switch lst := args[0].(type) {
	case types.Collection:
		definitions = lst.Data()
	default:
		return nil, nil, fmt.Errorf("invalid let* environment definition")

	}

	for i := 0; i < len(definitions); i += 2 {
		if _, err := evalDef(newEnv, definitions[i:]...); err != nil {
			return nil, nil, err
		}
	}

	return args[1], newEnv, nil
}

func evalDo(e types.Env, args ...types.Base) (types.Base, error) {
	_, err := evalAST(types.NewList(args[:len(args)-1]...), e)
	if err != nil {
		return nil, err
	}
	return args[len(args)-1], nil
}

func evalIf(e types.Env, args ...types.Base) (types.Base, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("improperly formatted if statement")
	}
	if condition, err := evalBool(e, args[0]); err != nil {
		return nil, err
	} else if condition {
		return args[1], nil
	} else if len(args) > 2 {
		return args[2], nil
	}
	return nil, nil
}

func evalFn(closureEnv types.Env, args ...types.Base) (types.Base, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("improperly formatted fn* statement")
	}

	var params []types.Base
	switch tparams := args[0].(type) {
	case types.Collection:
		params = tparams.Data()
	default:
		return nil, fmt.Errorf("invalid fn* param declaration")
	}

	fn := func(e types.Env, arguments []types.Base) (types.Base, error) {
		newEnv, err := env.New(closureEnv, params, arguments)
		if err != nil {
			return nil, err
		}
		return Eval(args[1], newEnv)
	}

	return &types.ExtFunc{
		AST:    args[1],
		Params: params,
		Env:    closureEnv,
		Fn:     fn,
	}, nil
}

func evalBool(e types.Env, condition types.Base) (bool, error) {
	value, err := Eval(condition, e)
	if err != nil {
		return false, err
	}
	switch tVal := value.(type) {
	case float64:
		return tVal == 0, nil
	case bool:
		return tVal, nil
	case nil:
		return false, nil
	default:
		return true, nil
	}
}

func evalQuasiQuote(e types.Env, object types.Base) types.Base {
	pair, isp := isPair(object)
	if !isp {
		return types.NewList(types.Symbol("quote"), object)
	}

	if sym, ok := pair[0].(types.Symbol); ok && sym == "unquote" {
		return pair[1]
	} else if nextPair, isp := isPair(pair[0]); isp {
		if sym, ok := nextPair[0].(types.Symbol); ok && sym == "splice-unquote" {
			return types.NewList(types.Symbol("concat"), nextPair[1], evalQuasiQuote(e, types.NewList(pair[1:]...)))
		}
	}
	return types.NewList(types.Symbol("cons"), evalQuasiQuote(e, pair[0]), evalQuasiQuote(e, types.NewList(pair[1:]...)))
}

func evalListForms(values []types.Base, env types.Env) ([]types.Base, error) {
	var err error
	forms := make([]types.Base, len(values))
	for i, form := range values {
		forms[i], err = Eval(form, env)
		if err != nil {
			return forms, err
		}
	}
	return forms, nil
}

func isPair(val types.Base) ([]types.Base, bool) {
	lst, isList := val.(types.Collection)
	if !isList {
		return []types.Base{}, false
	}
	data := lst.Data()
	return data, len(data) > 0
}
