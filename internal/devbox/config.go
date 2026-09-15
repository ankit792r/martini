package devbox

import (
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
		"echo 'Welcome to devbox!' > /dev/null",
		"alias dev='devbox'",
		fmt.Sprintf(`export PS1="(%s) [\$(pwd)] -> "`, projectName),
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
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
