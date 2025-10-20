package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}

func TestLoadPlaybookWithRoleReferences(t *testing.T) {
	tmpDir := t.TempDir()

	writeFile(t, filepath.Join(tmpDir, "roles", "vhosts", "common", "tasks", "main.yml"), `- name: Common task
  shell: echo common
`)

	writeFile(t, filepath.Join(tmpDir, "roles", "vhosts", "blackbox_exporter", "tasks", "main.yml"), `- name: Blackbox task
  shell: echo exporter
`)

	playbookContent := `- name: Deploy blackbox exporter
  hosts: test-hosts
  vars:
    hosts:
      - name: "www.svc.plus"
        path:
          - "/docs/"
  roles:
    - roles/vhosts/common/
    - role: vhosts/blackbox_exporter
  tasks:
    - name: Inline task
      shell: echo inline
`

	playbookPath := filepath.Join(tmpDir, "deploy.yml")
	writeFile(t, playbookPath, playbookContent)

	plays, err := LoadPlaybook(playbookPath)
	if err != nil {
		t.Fatalf("LoadPlaybook returned error: %v", err)
	}

	if len(plays) != 1 {
		t.Fatalf("expected 1 play, got %d", len(plays))
	}

	play := plays[0]
	if play.Hosts != "test-hosts" {
		t.Fatalf("unexpected hosts: %s", play.Hosts)
	}

	hostsVar, ok := play.Vars["hosts"].([]interface{})
	if !ok {
		t.Fatalf("expected hosts var to be a slice, got %T", play.Vars["hosts"])
	}
	if len(hostsVar) != 1 {
		t.Fatalf("expected 1 host entry, got %d", len(hostsVar))
	}

	if len(play.Tasks) != 3 {
		t.Fatalf("expected 3 tasks after expanding roles, got %d", len(play.Tasks))
	}

	if play.Tasks[0].Name != "Common task" {
		t.Fatalf("expected first task from common role, got %q", play.Tasks[0].Name)
	}
	if play.Tasks[1].Name != "Blackbox task" {
		t.Fatalf("expected second task from blackbox role, got %q", play.Tasks[1].Name)
	}
	if play.Tasks[2].Name != "Inline task" {
		t.Fatalf("expected inline task last, got %q", play.Tasks[2].Name)
	}
}

func TestLoadPlaybookWithLegacyModuleArgs(t *testing.T) {
	tmpDir := t.TempDir()

	writeFile(t, filepath.Join(tmpDir, "roles", "legacy", "tasks", "main.yml"), `- name: Render config
  template: src=blackbox.yml.j2 dest=/etc/blackbox.yml mode=0640
- name: Copy unit
  copy: src=blackbox.service dest=/etc/systemd/system/blackbox.service
- name: Remember vars
  set_fact:
    nested:
      key:
        - value
`)

	playbookContent := `- name: Legacy modules
  hosts: legacy-hosts
  roles:
    - roles/legacy
`

	playbookPath := filepath.Join(tmpDir, "deploy.yml")
	writeFile(t, playbookPath, playbookContent)

	plays, err := LoadPlaybook(playbookPath)
	if err != nil {
		t.Fatalf("LoadPlaybook returned error: %v", err)
	}
	if len(plays) != 1 {
		t.Fatalf("expected 1 play, got %d", len(plays))
	}

	tasks := plays[0].Tasks
	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}

	if tasks[0].Template == nil {
		t.Fatalf("expected template task to be parsed")
	}
	expectedTemplateSrc := filepath.Join(tmpDir, "roles", "legacy", "templates", "blackbox.yml.j2")
	if tasks[0].Template.Src != expectedTemplateSrc {
		t.Fatalf("unexpected template src: %s", tasks[0].Template.Src)
	}
	if tasks[0].Template.Dest != "/etc/blackbox.yml" {
		t.Fatalf("unexpected template dest: %s", tasks[0].Template.Dest)
	}
	if tasks[0].Template.Mode != "0640" {
		t.Fatalf("unexpected template mode: %s", tasks[0].Template.Mode)
	}

	if tasks[1].Copy == nil {
		t.Fatalf("expected copy task to be parsed")
	}
	expectedCopySrc := filepath.Join(tmpDir, "roles", "legacy", "files", "blackbox.service")
	if tasks[1].Copy.Src != expectedCopySrc {
		t.Fatalf("unexpected copy src: %s", tasks[1].Copy.Src)
	}
	if tasks[1].Copy.Dest != "/etc/systemd/system/blackbox.service" {
		t.Fatalf("unexpected copy dest: %s", tasks[1].Copy.Dest)
	}

	if tasks[2].SetFact == nil {
		t.Fatalf("expected set_fact to be parsed")
	}
	nested, ok := tasks[2].SetFact["nested"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected nested to be a map, got %T", tasks[2].SetFact["nested"])
	}
	values, ok := nested["key"].([]interface{})
	if !ok {
		t.Fatalf("expected key to be a slice, got %T", nested["key"])
	}
	if len(values) != 1 || values[0] != "value" {
		t.Fatalf("unexpected values: %#v", values)
	}
}
