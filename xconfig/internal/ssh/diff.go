package ssh

import "github.com/pmezard/go-difflib/difflib"

// Diff generates a unified diff string between before and after.
func Diff(before, after, dest string) string {
	ud := difflib.UnifiedDiff{
		A:        difflib.SplitLines(before),
		B:        difflib.SplitLines(after),
		FromFile: "before",
		ToFile:   "after: " + dest,
		Context:  3,
	}
	diff, _ := difflib.GetUnifiedDiffString(ud)
	return diff
}
