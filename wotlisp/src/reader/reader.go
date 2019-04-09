package reader

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/tanema/mal/wotlisp/src/types"
)

var (
	ErrUnderflow = errors.New("EOF underflow error: more input expected")

	tokensPattern = regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" + `~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" + `,;)]*)`)
	numberPattern = regexp.MustCompile(`^-?[0-9]+$`)
)

type Reader struct {
	tokens []string
}

func ReadString(in string) (types.Base, error) {
	reader := &Reader{tokens: tokenize(in)}
	return reader.form()
}

func (reader *Reader) next() (string, bool) {
	var token string
	if len(reader.tokens) == 0 {
		return token, false
	}
	token, reader.tokens = reader.tokens[0], reader.tokens[1:]
	return token, true
}

func (reader *Reader) peek() (string, bool) {
	if len(reader.tokens) > 0 {
		return reader.tokens[0], true
	}
	return "", false
}

func tokenize(in string) []string {
	results := []string{}
	for _, group := range tokensPattern.FindAllStringSubmatch(in, -1) {
		if (group[1] == "") || (group[1][0] == ';') {
			continue
		}
		results = append(results, group[1])
	}
	return results
}

func (reader *Reader) form() (types.Base, error) {
	token, hasNext := reader.peek()
	if !hasNext {
		return nil, ErrUnderflow
	}

	switch token {
	case `'`:
		return reader.modifier("quote")
	case "`":
		return reader.modifier("quasiquote")
	case `~`:
		return reader.modifier("unquote")
	case `~@`:
		return reader.modifier("splice-unquote")
	case `^`:
		return reader.meta()
	case `@`:
		return reader.modifier("deref")
	case ")":
		return nil, errors.New("unexpected ')'")
	case "(":
		return reader.list("(", ")")
	case "]":
		return nil, errors.New("unexpected ']'")
	case "[":
		return reader.vector()
	case "}":
		return nil, errors.New("unexpected '}'")
	case "{":
		return reader.hashMap()
	default:
		return reader.atom()
	}
}

func (reader *Reader) modifier(symbol string) (*types.List, error) {
	reader.next()
	form, err := reader.form()
	return types.NewList(types.Symbol(symbol), form), err
}

func (reader *Reader) meta() (*types.List, error) {
	reader.next()
	meta, err := reader.form()
	if err != nil {
		return nil, err
	}
	form, err := reader.form()
	return types.NewList(types.Symbol("with-meta"), form, meta), err
}

func (reader *Reader) list(start, end string) (*types.List, error) {
	list := &types.List{Forms: []types.Base{}}
	token, hasNext := reader.next()
	if !hasNext {
		return list, ErrUnderflow
	}
	if token != start {
		return list, fmt.Errorf("unexpected '%v'", token)
	}
	token, hasNext = reader.peek()
	for ; token != end && hasNext; token, hasNext = reader.peek() {
		form, err := reader.form()
		if err != nil {
			return list, err
		}
		list.Forms = append(list.Forms, form)
	}
	if token, hasNext := reader.next(); !hasNext {
		return list, ErrUnderflow
	} else if token != end {
		return list, fmt.Errorf("unexpected '%v'", token)
	}
	return list, nil
}

func (reader *Reader) vector() (*types.Vector, error) {
	list, err := reader.list("[", "]")
	return &types.Vector{Forms: list.Forms}, err
}

func (reader *Reader) hashMap() (*types.Hashmap, error) {
	list, err := reader.list("{", "}")
	if err != nil {
		return nil, err
	}
	return types.NewHashmap(list.Forms)
}

func (reader *Reader) atom() (types.Base, error) {
	token, hasNext := reader.next()
	if !hasNext {
		return nil, ErrUnderflow
	}

	if token == "nil" {
		return nil, nil
	} else if token == "true" {
		return true, nil
	} else if token == "false" {
		return false, nil
	} else if token[0] == ':' {
		return types.Keyword(token[1:]), nil
	} else if match := numberPattern.MatchString(token); match {
		var i int
		var e error
		if i, e = strconv.Atoi(token); e != nil {
			return nil, errors.New("number parse error")
		}
		return i, nil
	} else if token[0] == '"' {
		if token[len(token)-1] != '"' {
			return nil, errors.New("expected '\"', got EOF")
		}
		return token[1 : len(token)-1], nil
	}

	return types.Symbol(token), nil
}
