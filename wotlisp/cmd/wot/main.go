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
)

func main() {
	r := bufio.NewReader(os.Stdin)
	defaultEnv, _ := env.New(nil, nil, nil)
	for method, fn := range core.Namespace {
		defaultEnv.Set(method, fn)
	}
	rep("(def! not (fn* (a) (if a false true)))", defaultEnv)
	for {
		fmt.Print("user> ")
		text, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Println(rep(text, defaultEnv))
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
