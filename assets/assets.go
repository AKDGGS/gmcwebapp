package assets

import (
	"embed"
)

//go:embed img/* css/* pg/* tmpl/* ol/* js/*
var assets embed.FS

func ReadBytes(name string) []byte {
	f, _ := assets.ReadFile(name)
	return f
}

func ReadString(name string) string {
	f, _ := assets.ReadFile(name)
	return string(f)
}

func Initialize() error {
	if err := initStatic(); err != nil {
		return err
	}
	if err := initTemplates(); err != nil {
		return err
	}
	return nil
}
