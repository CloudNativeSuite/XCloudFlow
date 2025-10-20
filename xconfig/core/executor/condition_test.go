package executor

import (
	"testing"

	"xconfig/core/parser"
)

func TestEvaluateWhen(t *testing.T) {
	vars := map[string]interface{}{"flag": true, "count": 1, "empty": "", "number": 0}

	if !EvaluateWhen(parser.When{}, vars) {
		t.Fatalf("expected empty condition to evaluate true")
	}

	single := parser.When{Expressions: []string{"flag"}}
	if !EvaluateWhen(single, vars) {
		t.Fatalf("expected flag to evaluate true")
	}

	multi := parser.When{Expressions: []string{"flag", "count"}}
	if !EvaluateWhen(multi, vars) {
		t.Fatalf("expected both expressions to evaluate true")
	}

	falseExpr := parser.When{Expressions: []string{"empty"}}
	if EvaluateWhen(falseExpr, vars) {
		t.Fatalf("expected empty string to evaluate false")
	}

	zeroExpr := parser.When{Expressions: []string{"number"}}
	if EvaluateWhen(zeroExpr, vars) {
		t.Fatalf("expected zero to evaluate false")
	}
}
