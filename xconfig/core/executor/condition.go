package executor

import "xconfig/core/parser"

// EvaluateWhen returns true if the given when clause evaluates to true.
// Currently this supports basic variable lookup: if each expression corresponds
// to a variable name and that variable's value is truthy (not empty, "false" or
// "0"), the condition is true. An empty when clause evaluates to true.
func EvaluateWhen(when parser.When, vars map[string]interface{}) bool {
	if when.IsEmpty() {
		return true
	}
	for _, expr := range when.Expressions {
		if !evaluateExpression(expr, vars) {
			return false
		}
	}
	return true
}

func evaluateExpression(expr string, vars map[string]interface{}) bool {
	if expr == "" {
		return true
	}
	if val, ok := vars[expr]; ok {
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
	return expr == "true"
}
