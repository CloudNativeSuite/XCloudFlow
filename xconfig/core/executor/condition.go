package executor

import (
	"strings"

	"xconfig/core/parser"
)

// EvaluateWhen returns true if the expressions in the when clause evaluate to
// true. Currently this supports basic variable lookup: if the expression
// corresponds to a variable name and that variable's value is truthy (not empty,
// "false" or "0"), the condition is true. All expressions must evaluate to true
// for the task to run.
func EvaluateWhen(cond parser.When, vars map[string]interface{}) bool {
	if cond.IsEmpty() {
		return true
	}
	for _, expr := range cond.Expressions {
		if !evaluateExpression(expr, vars) {
			return false
		}
	}
	return true
}

func evaluateExpression(expr string, vars map[string]interface{}) bool {
	key := strings.TrimSpace(expr)
	if key == "" {
		return true
	}
	if val, ok := vars[key]; ok {
		switch v := val.(type) {
		case string:
			switch v {
			case "", "false", "0":
				return false
			default:
				return true
			}
		case bool:
			return v
		case int:
			return v != 0
		case int8:
			return v != 0
		case int16:
			return v != 0
		case int32:
			return v != 0
		case int64:
			return v != 0
		case uint:
			return v != 0
		case uint8:
			return v != 0
		case uint16:
			return v != 0
		case uint32:
			return v != 0
		case uint64:
			return v != 0
		case float32:
			return v != 0
		case float64:
			return v != 0
		default:
			return v != nil
		}
	}
	return key == "true"
}
