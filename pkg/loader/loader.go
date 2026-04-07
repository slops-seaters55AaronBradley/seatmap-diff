package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Load reads a YAML or JSON config file and returns the parsed data.
func Load(path string) (map[string]interface{}, error) {
	ext := filepath.Ext(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}

	var result map[string]interface{}
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &result)
	case ".json":
		err = json.Unmarshal(data, &result)
	default:
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return result, nil
}
