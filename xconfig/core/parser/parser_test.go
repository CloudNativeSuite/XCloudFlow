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
