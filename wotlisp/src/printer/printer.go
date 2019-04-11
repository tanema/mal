package printer

import (
	"fmt"
	"strings"

	"github.com/tanema/mal/wotlisp/src/types"
)

func List(forms []types.Base, pretty bool, pre, post, join string) string {
	strList := make([]string, len(forms))
	for i, e := range forms {
		strList[i] = Print(e, pretty)
	}
	return pre + strings.Join(strList, join) + post
}

func Print(object types.Base, pretty bool) string {
	switch tobj := object.(type) {
	case *types.Vector:
		return List(tobj.Forms, pretty, "[", "]", " ")
	case *types.List:
		return List(tobj.Forms, pretty, "(", ")", " ")
	case *types.Hashmap:
		return List(tobj.ToList(), pretty, "{", "}", " ")
	case types.Symbol:
		return string(tobj)
	case types.Keyword:
		return ":" + string(tobj)
	case types.Func:
		return "#<function>"
	case *types.ExtFunc:
		return "#<extfunction>"
	case *types.Atom:
		return "(atom " + Print(tobj.Val, pretty) + ")"
	case string:
		if pretty {
			tobj = strings.Replace(tobj, `\`, `\\`, -1)
			tobj = strings.Replace(tobj, `"`, `\"`, -1)
			tobj = strings.Replace(tobj, "\n", `\n`, -1)
			return `"` + tobj + `"`
		}
		return tobj
	case bool:
		if tobj {
			return "true"
		}
		return "false"
	case nil:
		return "nil"
	case float64:
		return fmt.Sprintf("%v", tobj)
	default:
		return "error formatting datatype"
	}
}
