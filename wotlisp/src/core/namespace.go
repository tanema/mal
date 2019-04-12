package core

import (
	"github.com/tanema/mal/wotlisp/src/env"
	"github.com/tanema/mal/wotlisp/src/runtime"
	"github.com/tanema/mal/wotlisp/src/types"
)

func DefaultNamespace() *env.Env {
	defaultEnv, _ := env.New(nil, nil, nil)
	for method, fn := range namespace {
		defaultEnv.Set(method, fn)
	}
	defaultEnv.Set("eval", eval(defaultEnv))
	evaluate(defaultEnv, "(def! not (fn* (a) (if a false true)))")
	evaluate(defaultEnv, `(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) ")")))))`)
	return defaultEnv
}

func evaluate(e *env.Env, source string) (types.Base, error) {
	ast, parseErr := readString(e, []types.Base{source})
	if parseErr != nil {
		return nil, parseErr
	}
	evalFn, err := e.Get("eval")
	if err != nil {
		return nil, err
	}
	return evalFn.(types.Func)(e, []types.Base{ast})
}

func eval(defaultEnv types.Env) types.Func {
	return func(e types.Env, a []types.Base) (types.Base, error) {
		if len(a) < 1 {
			return nil, nil
		}
		return runtime.Eval(a[0], defaultEnv)
	}
}
