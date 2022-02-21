package config

import (
	yaml "gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	ListenAddress string          `yaml:"listen_address"`
	DatabaseURL   string          `yaml:"database_url"`
	BasePath      string          `yaml:"base_path"`
	FileStore     FileStoreConfig `yaml:"file_store"`
	Auths         []AuthConfig    `yaml:"authentication"`
}

type FileStoreConfig struct {
	Type  string                 `yaml:"type"`
	Attrs map[string]interface{} `yaml:",inline"`
}

type AuthConfig struct {
	Type  string                 `yaml:"type"`
	Attrs map[string]interface{} `yaml:",inline"`
}

func New() *Config {
	cfg := &Config{
		ListenAddress: "127.0.0.1:8080",
		DatabaseURL:   "postgres://localhost",
		BasePath:      "/",
	}
	if cwd, err := os.Getwd(); err == nil {
		cfg.FileStore.Type = "dir"
		cfg.FileStore.Attrs = map[string]interface{}{
			"path": cwd,
		}
	}
	return cfg
}

func Load(cpath string) (*Config, error) {
	if cpath == "" {
		searchpaths := []string{"/etc/gmc.yaml"}
		if cdir, err := os.UserConfigDir(); err == nil {
			searchpaths = append(
				searchpaths, cdir+string(os.PathSeparator)+"gmc.yaml",
			)
		}
		searchpaths = append(searchpaths, "gmc.yaml")

		for _, sp := range searchpaths {
			if _, err := os.Stat(sp); err == nil {
				cpath = sp
				break
			}
		}
	}

	cfg := New()
	if cpath != "" {
		f, err := os.Open(cpath)
		if err != nil {
			return nil, err
		}

		cdata, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}

		err = yaml.UnmarshalStrict(cdata, &cfg)
		if err != nil {
			return nil, err
		}
	}
	return cfg, nil
}
