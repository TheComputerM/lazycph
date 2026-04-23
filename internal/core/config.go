package core

import (
	"encoding/json"
	"maps"
	"os"
	"path/filepath"
)

// default engines
var Engines = map[string]Engine{
	".c": {
		Mode:    "compile",
		Command: []string{"cc", "{file}", "-o", "{temp}", "-std=c17"},
	},
	".cpp": {
		Mode:    "compile",
		Command: []string{"c++", "{file}", "-o", "{temp}", "-std=c++17"},
	},
	".go": {
		Mode:    "compile",
		Command: []string{"go", "build", "-o", "{temp}", "{file}"},
	},
	".py": {
		Mode:    "interpret",
		Command: []string{"python3", "{file}"},
	},
	".rs": {
		Mode:    "compile",
		Command: []string{"rustc", "{file}", "-o", "{temp}"},
	},
	".zig": {
		Mode:    "compile",
		Command: []string{"zig", "build-exe", "{file}", "-femit-bin={temp}"},
	},
}

func init() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return
	}

	data, err := os.ReadFile(filepath.Join(configDir, "lazycph.json"))
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
