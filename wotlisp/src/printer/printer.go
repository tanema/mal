package printer

import (
	"fmt"
	"strings"

	"github.com/tanema/mal/wotlisp/src/types"
)

func formatlist(forms []types.Base, pre, post, join string) string {
	strList := make([]string, len(forms))
	for i, e := range forms {
		strList[i] = Print(e)
	}
	return pre + strings.Join(strList, join) + post
}

func Print(object types.Base) string {
	switch tobj := object.(type) {
	case *types.List:
		return formatlist(tobj.Forms, "(", ")", " ")
	case *types.Vector:
		return formatlist(tobj.Forms, "[", "]", " ")
	case *types.Hashmap:
		return formatlist(tobj.ToList(), "{", "}", " ")
	case types.Symbol:
		return string(tobj)
	case types.Keyword:
		return ":" + string(tobj)
	case string:
		return `"` + tobj + `"`
	case bool:
		if tobj {
			return "true"
		}
		return "false"
	case nil:
		return "nil"
	case int, int32, int64, uint, uint32, uint64, float32, float64:
		return fmt.Sprintf("%v", tobj)
	default:
		return "error formatting datatype"
	}
}
