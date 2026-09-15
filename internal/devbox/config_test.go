package devbox

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildConfig(t *testing.T) {
	env := map[string]string{"PATH": "$PATH:$HOME/.pub-cache/bin", "FOO": "bar"}
	scripts := map[string][]string{"build": {"flutter build"}}

	cfg := buildConfig("myapp", env, scripts)

	if cfg.Schema != schemaURL {
		t.Fatalf("unexpected schema: %s", cfg.Schema)
	}
	if len(cfg.Packages) != 1 || cfg.Packages[0] != "flutter@3.47.0-sdk-links" {
		t.Fatalf("unexpected packages: %v", cfg.Packages)
	}
	if cfg.Env["FOO"] != "bar" {
		t.Fatalf("unexpected env: %v", cfg.Env)
	}

	ps1 := cfg.Shell.InitHook[len(cfg.Shell.InitHook)-1]
	wantPS1 := `export PS1="(myapp) [\$(pwd)] -> "`
	if ps1 != wantPS1 {
		t.Fatalf("unexpected PS1 hook: %q", ps1)
	}
	if cfg.Shell.Scripts["build"][0] != "flutter build" {
		t.Fatalf("unexpected scripts: %v", cfg.Shell.Scripts)
	}
}

func TestBuildConfigDefaultScript(t *testing.T) {
	cfg := buildConfig("demo", defaultEnv(), nil)

	if _, ok := cfg.Shell.Scripts["test"]; !ok {
		t.Fatalf("expected default test script, got: %v", cfg.Shell.Scripts)
	}
}

func TestWriteConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "devbox.json")

	cfg := buildConfig("demo", defaultEnv(), nil)
	if err := writeConfig(path, cfg); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var decoded Config
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded.Schema != schemaURL {
		t.Fatalf("unexpected schema in file: %s", decoded.Schema)
	}
}
