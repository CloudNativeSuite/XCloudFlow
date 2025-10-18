package executor

// EvaluateWhen returns true if the given expression evaluates to true.
// Currently this supports basic variable lookup: if the expression corresponds
// to a variable name and that variable's value is truthy (not empty, "false" or
// "0"), the condition is true. Empty expression evaluates to true.
func EvaluateWhen(expr string, vars map[string]string) bool {
	if expr == "" {
		return true
	}
	if val, ok := vars[expr]; ok {
		switch val {
		case "", "false", "0":
			return false
		default:
			return true
		}
	}
	return expr == "true"
}
