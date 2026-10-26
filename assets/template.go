package assets

import (
	"html/template"
	"io"
	"reflect"
	"strings"
	"sync"
	"time"
)

type templateEntry struct {
	ModTime  time.Time
	Template *template.Template
}

var templateLock sync.RWMutex
var templateCache map[string]*templateEntry = make(map[string]*templateEntry)

func ExecuteTemplate(name string, wr io.Writer, data interface{}) error {
	s, err := Stat(name)
	if err != nil {
		return err
	}

	templateLock.RLock()
	curEntry, ok := templateCache[name]
	templateLock.RUnlock()

	if !ok || s.ModTime().After(curEntry.ModTime) {
		b, err := ReadBytes(name)
		if err != nil {
			return err
		}

		tmpl := template.New(name)
		tmpl.Funcs(template.FuncMap{
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
			"lower": strings.ToLower,
			"deref": func(v any) any {
				if v == nil {
					return ""
				}
				r := reflect.ValueOf(v)
				if r.Kind() != reflect.Ptr {
					return v
				}
				if r.IsNil() {
					return reflect.Zero(r.Type().Elem()).Interface()
				}
				return r.Elem().Interface()
			},
		})
		if _, err := tmpl.Parse(string(b)); err != nil {
			return err
		}

		curEntry = &templateEntry{ModTime: s.ModTime(), Template: tmpl}

		templateLock.Lock()
		templateCache[name] = curEntry
		templateLock.Unlock()
	}
	return curEntry.Template.Execute(wr, data)
}
