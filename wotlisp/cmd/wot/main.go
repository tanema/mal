package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tanema/mal/wotlisp/src/env"
	"github.com/tanema/mal/wotlisp/src/printer"
	"github.com/tanema/mal/wotlisp/src/reader"
	"github.com/tanema/mal/wotlisp/src/runtime"
	"github.com/tanema/mal/wotlisp/src/std"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	defaultEnv := env.New(nil)
	std.Define(defaultEnv)
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
	return printer.Print(val)
}
