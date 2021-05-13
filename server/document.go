package server

import (
	"net/http"
	"sync"

	"github.com/rakyll/statik/fs"
)

var once sync.Once
var document http.FileSystem

func defaultDocument() http.FileSystem {
	once.Do(func() {
		var e error
		document, e = fs.NewWithNamespace(`document`)
		if e != nil {
			panic(e)
		}
	})
	return document
}
