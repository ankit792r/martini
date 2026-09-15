package devbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

const schemaURL = "https://raw.githubusercontent.com/jetify-com/devbox/0.18.0/.schema/devbox.schema.json"

type Config struct {
	Schema   string            `json:"$schema"`
	Packages []string          `json:"packages"`
	Env      map[string]string `json:"env"`
	Shell    Shell             `json:"shell"`
}

type Shell struct {
	InitHook []string            `json:"init_hook"`
	Scripts  map[string][]string `json:"scripts"`
}

func buildConfig(projectName string, env map[string]string, scripts map[string][]string) Config {
	initHook := []string{
		"alias dev='devbox'",
		fmt.Sprintf(`export PS1="(%s) [\W] -> "`, projectName),
	}

	if len(scripts) == 0 {
		scripts = map[string][]string{
			"test": {`echo "Error: no test specified" && exit 1`},
		}
	}

	return Config{
		Schema:   schemaURL,
		Packages: []string{},
		Env:      env,
		Shell: Shell{
			InitHook: initHook,
			Scripts:  scripts,
		},
	}
}

func writeConfig(path string, cfg Config) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cfg); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}
