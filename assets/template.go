package assets

import (
	"html/template"
	"io"
)

// Embedded templates can be safely cached
var embeddedTemplates *template.Template
var stdFuncs template.FuncMap

func initTemplates() error {
	stdFuncs = template.FuncMap{
		"rawjs": func(in interface{}) interface{} {
			switch v := in.(type) {
			case string:
				return template.JS(v)
			case nil:
				return nil
			default:
				return in
			}
		},
		"rawhtml": func(in interface{}) interface{} {
			switch v := in.(type) {
			case string:
				return template.HTML(v)
			case nil:
				return nil
			default:
				return in
			}
		},
	}

	embeddedTemplates = template.New("templates").Funcs(stdFuncs)
	_, err := embeddedTemplates.ParseFS(assets, "tmpl/*.html")
	if err != nil {
		return err
	}
	return nil
}

func LookupTemplate(name string) *template.Template {
	return embeddedTemplates.Lookup(name)
}

func ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return embeddedTemplates.ExecuteTemplate(wr, name, data)
}
