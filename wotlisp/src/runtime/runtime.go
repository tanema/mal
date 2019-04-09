package runtime

import (
	"fmt"

	"github.com/tanema/mal/wotlisp/src/env"
	"github.com/tanema/mal/wotlisp/src/types"
)

func Eval(object types.Base, e *env.Env) (types.Base, error) {
	switch tobject := object.(type) {
	case *types.List:
		if len(tobject.Forms) == 0 {
			return tobject, nil
		}
		sym, _ := tobject.Forms[0].(types.Symbol)
		switch sym {
		case "def!":
			return evalDef(e, tobject.Forms[1:]...)
		case "let*":
			return evalLet(e, tobject.Forms[1:]...)
		default:
			return evalFnCall(e, tobject)
		}
	default:
		return evalAST(tobject, e)
	}
}

func evalAST(ast types.Base, env *env.Env) (types.Base, error) {
	switch tobject := ast.(type) {
	case types.Symbol:
		symVal, err := env.Get(tobject)
		if err != nil {
			return nil, err
		}
		return symVal, nil
	case *types.List:
		lst, err := evalListForms(tobject.Forms, env)
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

func evalDef(e *env.Env, args ...types.Base) (types.Base, error) {
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

func evalLet(e *env.Env, args ...types.Base) (types.Base, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("not enough arguments for  let* call")
	}
	newEnv := env.New(e)

	var definitions []types.Base
	switch lst := args[0].(type) {
	case *types.List:
		definitions = lst.Forms
	case *types.Vector:
		definitions = lst.Forms
	default:
		return nil, fmt.Errorf("invalid let* environment definition")

	}

	for i := 0; i < len(definitions); i += 2 {
		if _, err := evalDef(newEnv, definitions[i:]...); err != nil {
			return nil, err
		}
	}

	return Eval(args[1], newEnv)
}

func evalFnCall(e *env.Env, list *types.List) (types.Base, error) {
	lst, err := evalAST(list, e)
	if err != nil {
		return nil, err
	}
	list = lst.(*types.List)
	fn, ok := list.Forms[0].(func(*env.Env, []types.Base) (types.Base, error))
	if !ok {
		return nil, fmt.Errorf("attempt to call non-function %v", list.Forms[0])
	}
	return fn(e, list.Forms[1:])
}

func evalListForms(values []types.Base, env *env.Env) ([]types.Base, error) {
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
