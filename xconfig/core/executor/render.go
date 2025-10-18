package executor

import (
	"bytes"
	"text/template"
)

// RenderString renders a template string with the provided variables.
func RenderString(tmplStr string, vars map[string]string) (string, error) {
	t, err := template.New("tmpl").Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, vars); err != nil {
		return "", err
	}
	return buf.String(), nil
}
