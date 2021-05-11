package configure

import (
	"encoding/json"

	"github.com/google/go-jsonnet"
)

var defaultConfigure Configure

func DefaultConfigure() *Configure {
	return &defaultConfigure
}

type Configure struct {
	Server   Server
	Manager  Manager
	Provider Provider
}

func Load(filename string) (e error) {
	vm := jsonnet.MakeVM()
	vm.Importer(&jsonnet.FileImporter{})
	jsonStr, e := vm.EvaluateFile(filename)
	if e != nil {
		return
	}
	cnf := DefaultConfigure()
	e = json.Unmarshal([]byte(jsonStr), cnf)
	return
}
