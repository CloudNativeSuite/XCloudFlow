package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

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

	fi, err := os.Stat(configDir)
	if err != nil || !fi.IsDir() {
		return nil, fmt.Errorf("\u274c \u914d\u7f6e\u76ee\u5f55\u4e0d\u5b58\u5728: %s", configDir)
	}

	patterns := []string{"*.yaml", "*.yml"}
	var files []string
	for _, p := range patterns {
		list, _ := filepath.Glob(filepath.Join(configDir, p))
		files = append(files, list...)
	}
	sort.Strings(files)
	if len(files) == 0 {
		return nil, fmt.Errorf("\u26a0\ufe0f \u672a\u627e\u5230\u4efb\u4f55 YAML \u914d\u7f6e\u6587\u4ef6\u4e8e: %s", configDir)
	}

	merged := make(map[string]interface{})
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
	merged["__config_path__"] = configDir
	return merged, nil
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
