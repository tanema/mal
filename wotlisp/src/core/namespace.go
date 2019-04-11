package core

import (
	"github.com/tanema/mal/wotlisp/src/env"
	"github.com/tanema/mal/wotlisp/src/types"
)

func DefaultNamespace() *env.Env {
	defaultEnv, _ := env.New(nil, nil, nil)
	for method, fn := range namespace {
		defaultEnv.Set(method, fn)
	}
	defaultEnv.Set("eval", eval(defaultEnv))
	Eval(defaultEnv, "(def! not (fn* (a) (if a false true)))")
	Eval(defaultEnv, `(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) ")")))))`)
	return defaultEnv
}

func Eval(e *env.Env, source string, argv ...string) (types.Base, error) {
	ast, parseErr := readString(e, []types.Base{source})
	if parseErr != nil {
		return nil, parseErr
	}
	evalFn, err := e.Get("eval")
	if err != nil {
		return nil, err
	}
	targv := make([]types.Base, len(argv))
	for i, arg := range argv {
		targv[i] = types.Base(arg)
	}
	e.Set("*ARGV*", types.NewList(targv...))
	return evalFn.(types.Func)(e, []types.Base{ast})
}
