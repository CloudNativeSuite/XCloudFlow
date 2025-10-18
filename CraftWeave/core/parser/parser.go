// core/parser/parser.go
package parser

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Template struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}

type Copy struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
	Mode string `yaml:"mode,omitempty"`
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
	Name     string            `yaml:"name"`
	When     string            `yaml:"when,omitempty"`
	Shell    string            `yaml:"shell,omitempty"`
	Script   string            `yaml:"script,omitempty"`
	Template *Template         `yaml:"template,omitempty"`
	Command  string            `yaml:"command,omitempty"`
	Copy     *Copy             `yaml:"copy,omitempty"`
	Stat     *Stat             `yaml:"stat,omitempty"`
	Apt      *PackageAction    `yaml:"apt,omitempty"`
	Yum      *PackageAction    `yaml:"yum,omitempty"`
	Systemd  *ServiceAction    `yaml:"systemd,omitempty"`
	Service  *ServiceAction    `yaml:"service,omitempty"`
	Setup    bool              `yaml:"setup,omitempty"`
	SetFact  map[string]string `yaml:"set_fact,omitempty"`
	Fail     *MessageAction    `yaml:"fail,omitempty"`
	Debug    *MessageAction    `yaml:"debug,omitempty"`
	Vultr    *VultrInstance    `yaml:"vultr,omitempty"`
	Register string            `yaml:"register,omitempty"`
}

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

type Play struct {
	Name  string            `yaml:"name"`
	Hosts string            `yaml:"hosts"`
	Vars  map[string]string `yaml:"vars,omitempty"`
	Roles []struct {
		Role string `yaml:"role"`
	} `yaml:"roles,omitempty"`
	Tasks []Task `yaml:"tasks,omitempty"`
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
			ts, err := loadRoleTasks(base, r.Role)
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
	roleDir := filepath.Join(base, "roles", name)
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
