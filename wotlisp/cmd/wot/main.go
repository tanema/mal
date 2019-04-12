package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tanema/mal/wotlisp/src/core"
	"github.com/tanema/mal/wotlisp/src/env"
	"github.com/tanema/mal/wotlisp/src/printer"
	"github.com/tanema/mal/wotlisp/src/reader"
	"github.com/tanema/mal/wotlisp/src/runtime"
	"github.com/tanema/mal/wotlisp/src/types"
)

func main() {
	defaultEnv := core.DefaultNamespace()
	if len(os.Args) > 1 {
		runFile(defaultEnv, os.Args[1], os.Args[2:]...)
	} else {
		runREPL(defaultEnv)
	}
}

func runFile(e *env.Env, path string, argv ...string) {
	targv := make([]types.Base, len(argv))
	for i, arg := range argv {
		targv[i] = types.Base(arg)
	}
	e.Set("*ARGV*", types.NewList(targv...))
	fmt.Println(rep(`(load-file "`+os.Args[1]+`")`, e))
}

func runREPL(env *env.Env) error {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("user> ")
		text, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		fmt.Println(rep(text, env))
	}
}

func rep(in string, e *env.Env) string {
	ast, parseErr := reader.ReadString(in)
	if parseErr != nil {
		return parseErr.Error()
	}
	val, evalErr := runtime.Eval(ast, e)
	if evalErr != nil {
		return evalErr.Error()
	}
	return printer.Print(val, true)
}
