package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	ListenAddress string          `yaml:"listen_address"`
	DatabaseURL   string          `yaml:"database_url"`
	BasePath      string          `yaml:"base_path"`
	FileStore     FileStoreConfig `yaml:"file_store"`
	SessionKey    string          `yaml:"session_key"`
	keybytes      []byte          `yaml:"-"`
	MaxAge        int             `yaml:"session_max_age"`
	Auths         []AuthConfig    `yaml:"authentication"`
}

func (c *Config) SessionKeyBytes() []byte {
	return c.keybytes
}

type FileStoreConfig struct {
	Type  string                 `yaml:"type"`
	Attrs map[string]interface{} `yaml:",inline"`
}

type AuthConfig struct {
	Type  string                 `yaml:"type"`
	Attrs map[string]interface{} `yaml:",inline"`
}

func New() (*Config, error) {
	cfg := &Config{
		ListenAddress: "127.0.0.1:8080",
		DatabaseURL:   "postgres://localhost",
		BasePath:      "/",
		MaxAge:        86400, // Default of 24 hours
	}

	cfg.keybytes = make([]byte, 32)
	_, err := rand.Read(cfg.keybytes)
	if err != nil {
		return nil, err
	}

	if cwd, err := os.Getwd(); err == nil {
		cfg.FileStore.Type = "dir"
		cfg.FileStore.Attrs = map[string]interface{}{
			"path": cwd,
		}
	}

	return cfg, nil
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

	cfg, err := New()
	if err != nil {
		return nil, err
	}

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

		if cfg.SessionKey != "" {
			b, err := hex.DecodeString(cfg.SessionKey)
			if err != nil || len(b) != 32 {
				return nil, fmt.Errorf(
					"session_key must be exactly 64 hexidecimal characters",
				)
			}

			cfg.keybytes = b
		}
	}

	return cfg, nil
}
