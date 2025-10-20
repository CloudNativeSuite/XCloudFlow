// core/parser/parser.go
package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Template struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
	Mode string `yaml:"mode,omitempty"`
}

// UnmarshalYAML allows templates to be specified using Ansible's inline
// "key=value" argument syntax or a regular mapping.
func (t *Template) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		args := parseModuleArgs(value.Value)
		if args == nil {
			return fmt.Errorf("invalid template arguments: %q", value.Value)
		}
		t.Src = args["src"]
		t.Dest = args["dest"]
		t.Mode = args["mode"]
		return nil
	case yaml.MappingNode:
		type templateAlias Template
		var tmp templateAlias
		if err := value.Decode(&tmp); err != nil {
			return err
		}
		*t = Template(tmp)
		return nil
	default:
		return fmt.Errorf("unsupported template format: %v", value.Kind)
	}
}

type Copy struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
	Mode string `yaml:"mode,omitempty"`
}

// UnmarshalYAML allows copy tasks to be specified either as a mapping or using
// inline "key=value" arguments.
func (c *Copy) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		args := parseModuleArgs(value.Value)
		if args == nil {
			return fmt.Errorf("invalid copy arguments: %q", value.Value)
		}
		c.Src = args["src"]
		c.Dest = args["dest"]
		c.Mode = args["mode"]
		return nil
	case yaml.MappingNode:
		type copyAlias Copy
		var tmp copyAlias
		if err := value.Decode(&tmp); err != nil {
			return err
		}
		*c = Copy(tmp)
		return nil
	default:
		return fmt.Errorf("unsupported copy format: %v", value.Kind)
	}
}

type Stat struct {
	Path string `yaml:"path"`
}

type PackageAction struct {
	Name  string `yaml:"name,omitempty"`
	Deb   string `yaml:"deb,omitempty"`
	State string `yaml:"state,omitempty"`
}

type ServiceAction struct {
	Name    string `yaml:"name"`
	State   string `yaml:"state"`
	Enabled bool   `yaml:"enabled,omitempty"`
}

type MessageAction struct {
	Msg string `yaml:"msg"`
}

// VultrInstance defines parameters to create a Vultr cloud instance.
type VultrInstance struct {
	APIKey string `yaml:"api_key,omitempty"`
	Region string `yaml:"region"`
	Plan   string `yaml:"plan"`
	OsID   int    `yaml:"os_id"`
	Label  string `yaml:"label,omitempty"`
}

type Task struct {
	Name     string                 `yaml:"name"`
	When     When                   `yaml:"when,omitempty"`
	Shell    string                 `yaml:"shell,omitempty"`
	Script   string                 `yaml:"script,omitempty"`
	Template *Template              `yaml:"template,omitempty"`
	Command  string                 `yaml:"command,omitempty"`
	Copy     *Copy                  `yaml:"copy,omitempty"`
	Stat     *Stat                  `yaml:"stat,omitempty"`
	Apt      *PackageAction         `yaml:"apt,omitempty"`
	Yum      *PackageAction         `yaml:"yum,omitempty"`
	Systemd  *ServiceAction         `yaml:"systemd,omitempty"`
	Service  *ServiceAction         `yaml:"service,omitempty"`
	Setup    bool                   `yaml:"setup,omitempty"`
	SetFact  map[string]interface{} `yaml:"set_fact,omitempty"`
	Fail     *MessageAction         `yaml:"fail,omitempty"`
	Debug    *MessageAction         `yaml:"debug,omitempty"`
	Vultr    *VultrInstance         `yaml:"vultr,omitempty"`
	Register string                 `yaml:"register,omitempty"`
}

// When represents one or more conditional expressions attached to a task.
type When struct {
	Expressions []string
}

// UnmarshalYAML accepts either a single string or a sequence of strings for the
// `when` clause, matching Ansible's syntax.
func (w *When) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case 0:
		return nil
	case yaml.ScalarNode:
		var expr string
		if err := value.Decode(&expr); err != nil {
			return err
		}
		expr = strings.TrimSpace(expr)
		if expr != "" {
			w.Expressions = []string{expr}
		}
		return nil
	case yaml.SequenceNode:
		var exprs []string
		for _, node := range value.Content {
			var expr string
			if err := node.Decode(&expr); err != nil {
				return err
			}
			expr = strings.TrimSpace(expr)
			if expr != "" {
				exprs = append(exprs, expr)
			}
		}
		w.Expressions = exprs
		return nil
	default:
		return fmt.Errorf("unsupported when format: %v", value.Kind)
	}
}

// IsEmpty returns true when no expressions are defined.
func (w When) IsEmpty() bool { return len(w.Expressions) == 0 }

// Type returns the module name associated with this task.
func (t Task) Type() string {
	switch {
	case t.Shell != "":
		return "shell"
	case t.Command != "":
		return "command"
	case t.Script != "":
		return "script"
	case t.Template != nil:
		return "template"
	case t.Copy != nil:
		return "copy"
	case t.Stat != nil:
		return "stat"
	case t.Apt != nil:
		return "apt"
	case t.Yum != nil:
		return "yum"
	case t.Systemd != nil:
		return "systemd"
	case t.Service != nil:
		return "service"
	case t.Setup:
		return "setup"
	case len(t.SetFact) > 0:
		return "set_fact"
	case t.Fail != nil:
		return "fail"
	case t.Debug != nil:
		return "debug"
	case t.Vultr != nil:
		return "vultr_instance"
	default:
		return ""
	}
}

type RoleRef struct {
	Name string
}

// UnmarshalYAML allows RoleRef to be specified either as a string or as a
// mapping with a "role" key, mirroring Ansible's playbook syntax.
func (r *RoleRef) UnmarshalYAML(value *yaml.Node) error {
	switch value.Kind {
	case yaml.ScalarNode:
		var name string
		if err := value.Decode(&name); err != nil {
			return err
		}
		r.Name = name
		return nil
	case yaml.MappingNode:
		var tmp struct {
			Role string `yaml:"role"`
		}
		if err := value.Decode(&tmp); err != nil {
			return err
		}
		if tmp.Role == "" {
			return fmt.Errorf("role mapping missing 'role' key")
		}
		r.Name = tmp.Role
		return nil
	default:
		return fmt.Errorf("unsupported role format: %v", value.Kind)
	}
}

type Play struct {
	Name  string                 `yaml:"name"`
	Hosts string                 `yaml:"hosts"`
	Vars  map[string]interface{} `yaml:"vars,omitempty"`
	Roles []RoleRef              `yaml:"roles,omitempty"`
	Tasks []Task                 `yaml:"tasks,omitempty"`
}

// LoadPlaybook parses the given playbook YAML and expands any referenced roles.
func LoadPlaybook(path string) ([]Play, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var plays []Play
	if err := yaml.Unmarshal(data, &plays); err != nil {
		return nil, err
	}

	base := filepath.Dir(path)
	for i := range plays {
		var allTasks []Task
		for _, r := range plays[i].Roles {
			ts, err := loadRoleTasks(base, r.Name)
			if err != nil {
				return nil, err
			}
			allTasks = append(allTasks, ts...)
		}
		allTasks = append(allTasks, plays[i].Tasks...)
		plays[i].Tasks = allTasks
	}

	return plays, nil
}

func loadRoleTasks(base, name string) ([]Task, error) {
	cleanName := strings.TrimSuffix(name, string(filepath.Separator))
	cleanName = filepath.Clean(cleanName)

	candidates := []string{
		filepath.Join(base, cleanName),
		filepath.Join(base, "roles", cleanName),
	}

	var roleDir string
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			roleDir = candidate
			break
		}
	}

	if roleDir == "" {
		return nil, fmt.Errorf("role '%s' not found", name)
	}

	dir := filepath.Join(roleDir, "tasks")
	path := filepath.Join(dir, "main.yaml")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = filepath.Join(dir, "main.yml")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var tasks []Task
	if err := yaml.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	for i := range tasks {
		if tasks[i].Script != "" && !filepath.IsAbs(tasks[i].Script) {
			tasks[i].Script = filepath.Join(roleDir, "scripts", tasks[i].Script)
		}
		if tasks[i].Template != nil && tasks[i].Template.Src != "" && !filepath.IsAbs(tasks[i].Template.Src) {
			tasks[i].Template.Src = filepath.Join(roleDir, "templates", tasks[i].Template.Src)
		}
		if tasks[i].Copy != nil && tasks[i].Copy.Src != "" && !filepath.IsAbs(tasks[i].Copy.Src) {
			tasks[i].Copy.Src = filepath.Join(roleDir, "files", tasks[i].Copy.Src)
		}
	}
	return tasks, nil
}

func parseModuleArgs(input string) map[string]string {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	type stateType int
	const (
		stateKey stateType = iota
		stateValue
	)

	result := make(map[string]string)
	var keyBuilder strings.Builder
	var valueBuilder strings.Builder
	state := stateKey
	var quote rune

	flushPair := func() {
		key := strings.TrimSpace(keyBuilder.String())
		if key == "" {
			keyBuilder.Reset()
			valueBuilder.Reset()
			return
		}
		value := strings.TrimSpace(valueBuilder.String())
		result[key] = value
		keyBuilder.Reset()
		valueBuilder.Reset()
	}

	for _, r := range input {
		switch state {
		case stateKey:
			switch {
			case r == '=':
				state = stateValue
			case r == ' ' || r == '\t':
				// ignore whitespace between arguments
				if keyBuilder.Len() > 0 {
					// standalone key without value is invalid
					return nil
				}
			default:
				keyBuilder.WriteRune(r)
			}
		case stateValue:
			if quote != 0 {
				if r == quote {
					quote = 0
				} else {
					valueBuilder.WriteRune(r)
				}
				continue
			}
			switch r {
			case '\'', '"':
				quote = r
			case ' ', '\t':
				flushPair()
				state = stateKey
			default:
				valueBuilder.WriteRune(r)
			}
		}
	}

	if keyBuilder.Len() > 0 || valueBuilder.Len() > 0 {
		flushPair()
	}

	return result
}
