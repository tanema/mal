package runtime

import (
	"fmt"

	"github.com/tanema/mal/wotlisp/src/types"
)

func evalList(values []types.Base, env map[string]types.Base) ([]types.Base, error) {
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

func Eval(object types.Base, env map[string]types.Base) (types.Base, error) {
	switch tobject := object.(type) {
	case *types.List:
		if len(tobject.Forms) == 0 {
			return tobject, nil
		}
		lst, err := evalAST(tobject, env)
		if err != nil {
			return nil, err
		}
		list := lst.(*types.List)
		fn, ok := list.Forms[0].(func([]types.Base) (types.Base, error))
		if !ok {
			return nil, fmt.Errorf("attempt to call non-function %v", list.Forms[0])
		}
		return fn(list.Forms[1:])
	default:
		return evalAST(tobject, env)
	}
}

func evalAST(ast types.Base, env map[string]types.Base) (types.Base, error) {
	switch tobject := ast.(type) {
	case types.Symbol:
		symVal, found := env[string(tobject)]
		if !found {
			return nil, fmt.Errorf("undefined symbol %v", string(tobject))
		}
		return symVal, nil
	case *types.List:
		lst, err := evalList(tobject.Forms, env)
		return &types.List{Forms: lst}, err
	case *types.Vector:
		lst, err := evalList(tobject.Forms, env)
		return &types.Vector{Forms: lst}, err
	case *types.Hashmap:
		lst, err := evalList(tobject.ToList(), env)
		if err != nil {
			return nil, err
		}
		return types.NewHashmap(lst)
	default:
		return ast, nil
	}
}
