package assets

import (
	"html/template"
	"io"
	"sync"
	"time"
)

// Standard functions for templates
var stdFuncs template.FuncMap = template.FuncMap{
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
	"formattime": func(fmt string, in interface{}) interface{} {
		switch v := in.(type) {
		case time.Time:
			return v.Format(fmt)
		case nil:
			return nil
		default:
			return in
		}
	},
}

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

		t, err := template.New(name).Funcs(stdFuncs).Parse(string(b))
		if err != nil {
			return err
		}

		curEntry = &templateEntry{ModTime: s.ModTime(), Template: t}

		templateLock.Lock()
		templateCache[name] = curEntry
		templateLock.Unlock()
	}
	return curEntry.Template.Execute(wr, data)
}
