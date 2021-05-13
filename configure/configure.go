package configure

import (
	"encoding/json"

	"github.com/google/go-jsonnet"
	"github.com/powerpuffpenguin/sessionid_server/logger"
)

var defaultConfigure Configure

func DefaultConfigure() *Configure {
	return &defaultConfigure
}

type Configure struct {
	Auth     Auth
	Server   Server
	Manager  Manager
	Provider Provider
	Logger   logger.Options
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
