package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/powerpuffpenguin/sessionid_server/configure"
	"github.com/powerpuffpenguin/sessionid_server/version"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var (
		help, ver bool
		cnf       string
	)
	flag.BoolVar(&help, `help`, false, `display help`)
	flag.BoolVar(&ver, `version`, false, `display version`)
	flag.StringVar(&cnf, `conf`, `server.jsonnet`, `configure file`)

	flag.Parse()
	if help {
		flag.PrintDefaults()
	} else if ver {
		fmt.Println(version.Platform)
		fmt.Println(version.Version)
	} else {
		e := configure.Load(cnf)
		if e != nil {
			log.Fatalln(e)
		}
	}
}
