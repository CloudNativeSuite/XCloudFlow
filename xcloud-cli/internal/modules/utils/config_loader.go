package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func deepMerge(dst, src map[string]interface{}) map[string]interface{} {
	if dst == nil {
		dst = map[string]interface{}{}
	}
	for k, v := range src {
		if dv, ok := dst[k]; ok {
			switch dvTyped := dv.(type) {
			case map[string]interface{}:
				if sv, ok := v.(map[string]interface{}); ok {
					dst[k] = deepMerge(dvTyped, sv)
					continue
				}
			case []interface{}:
				if sv, ok := v.([]interface{}); ok {
					dst[k] = append(dvTyped, sv...)
					continue
				}
			}
		}
		dst[k] = v
	}
	return dst
}

func LoadMergedConfig(configDir string) (map[string]interface{}, error) {
	if configDir == "" {
		configDir = os.Getenv("CONFIG_PATH")
		if configDir == "" {
			configDir = "config"
		}
	}

	env := os.Getenv("STACK_ENV")
	if env == "" {
		env = "sit"
	}

	dirs, err := resolveConfigDirs(configDir, env)
	if err != nil {
		return nil, err
	}

	merged := make(map[string]interface{})
	for _, dir := range dirs {
		files, err := listYAMLFiles(dir)
		if err != nil {
			return nil, err
		}
		if len(files) == 0 {
			continue
		}
		for _, f := range files {
			data, err := os.ReadFile(f)
			if err != nil {
				return nil, err
			}
			var part map[string]interface{}
			if err := yaml.Unmarshal(data, &part); err != nil {
				return nil, err
			}
			merged = deepMerge(merged, part)
		}
	}
	if len(merged) == 0 {
		return nil, fmt.Errorf("\u26a0\ufe0f \u672a\u627e\u5230\u4efb\u4f55 YAML \u914d\u7f6e\u6587\u4ef6\u4e8e: %s", strings.Join(dirs, ","))
	}
	merged["__config_path__"] = strings.Join(dirs, ",")
	merged["__active_env__"] = env
	return merged, nil
}

func resolveConfigDirs(configDir, env string) ([]string, error) {
	fi, err := os.Stat(configDir)
	if err != nil {
		return nil, fmt.Errorf("\u274c \u914d\u7f6e\u76ee\u5f55\u4e0d\u5b58\u5728: %s", configDir)
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf("\u274c \u914d\u7f6e\u76ee\u5f55\u4e0d\u662f\u76ee\u5f55: %s", configDir)
	}

	envDir := filepath.Join(configDir, env)
	if info, err := os.Stat(envDir); err == nil && info.IsDir() {
		dirs := make([]string, 0, 3)
		if baseInfo, err := os.Stat(filepath.Join(configDir, "base")); err == nil && baseInfo.IsDir() {
			dirs = append(dirs, filepath.Join(configDir, "base"))
		}
		dirs = append(dirs, envDir)
		return dirs, nil
	}

	// configDir may already point to the environment directory
	dirs := []string{}
	parent := filepath.Dir(configDir)
	if baseInfo, err := os.Stat(filepath.Join(parent, "base")); err == nil && baseInfo.IsDir() {
		dirs = append(dirs, filepath.Join(parent, "base"))
	}
	dirs = append(dirs, configDir)
	return dirs, nil
}

func listYAMLFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() != "." && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml") {
			files = append(files, path)
		}
		return nil
	})
	sort.Strings(files)
	return files, err
}

// DecodeSection decodes a section of map into out.
func DecodeSection(m map[string]interface{}, key string, out interface{}) error {
	sec, ok := m[key]
	if !ok {
		return fmt.Errorf("missing section %s", key)
	}
	b, err := yaml.Marshal(sec)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, out)
}
