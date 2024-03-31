package cache

import (
	"strings"
)

const (
	PLAIN  = 0
	BROTLI = 1 << iota
	GZIP
)

var cmap Map = Map{entries: map[string]*Entry{}}

func Get(name string) *Entry {
	return cmap.Get(name)
}

func Put(name string, entry *Entry) {
	cmap.Put(name, entry)
}

func Remove(name string) {
	cmap.Remove(name)
}

func PurgeExpired() {
	cmap.PurgeExpired()
}

func parseEncoding(accept string) int {
	r := PLAIN
	for _, v := range strings.Split(accept, ",") {
		v = strings.ToLower(strings.TrimSpace(v))
		switch v {
		case "br":
			r = r | BROTLI
		case "gzip":
			r = r | GZIP
		}
	}
	return r
}
