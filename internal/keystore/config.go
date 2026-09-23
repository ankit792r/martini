package keystore

import "strings"

const DefaultKeyAlias = "upload"

type Config struct {
	Name          string
	StorePassword string
	KeyPassword   string
	Alias         string
	CommonName    string
}

func ConfigFromAnswers(answers []string) Config {
	cfg := Config{
		Name:       "upload-keystore",
		Alias:      DefaultKeyAlias,
		CommonName: "Android Upload",
	}

	if len(answers) > 0 && strings.TrimSpace(answers[0]) != "" {
		cfg.Name = strings.TrimSpace(strings.TrimSuffix(answers[0], ".jks"))
	}
	if len(answers) > 1 {
		cfg.StorePassword = strings.TrimSpace(answers[1])
	}
	if len(answers) > 2 && strings.TrimSpace(answers[2]) != "" {
		cfg.KeyPassword = strings.TrimSpace(answers[2])
	} else {
		cfg.KeyPassword = cfg.StorePassword
	}
	if len(answers) > 3 && strings.TrimSpace(answers[3]) != "" {
		cfg.Alias = strings.TrimSpace(answers[3])
	}
	if len(answers) > 4 && strings.TrimSpace(answers[4]) != "" {
		cfg.CommonName = strings.TrimSpace(answers[4])
	}

	return cfg
}

func (cfg Config) Normalized() Config {
	if strings.TrimSpace(cfg.KeyPassword) == "" {
		cfg.KeyPassword = cfg.StorePassword
	}
	if strings.TrimSpace(cfg.Alias) == "" {
		cfg.Alias = DefaultKeyAlias
	}
	if strings.TrimSpace(cfg.CommonName) == "" {
		cfg.CommonName = "Android Upload"
	}
	return cfg
}
