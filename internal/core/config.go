package core

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
)

// default engines
var Engines = map[string]Engine{
	"cpp": {
		Mode:    "compile",
		Command: []string{"c++", "{file}", "-o", "{temp}", "-std=c++17"},
	},
	"py": {
		Mode:    "interpret",
		Command: []string{"python3", "{file}"},
	},
}

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return
	}

	data, err := os.ReadFile(filepath.Join(configDir, "lazycph", "lazycph.json"))
	if err != nil {
		return
	}

	var cfg struct {
		Engines map[string]Engine `json:"engines"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return
	}

	maps.Copy(Engines, cfg.Engines)
}
