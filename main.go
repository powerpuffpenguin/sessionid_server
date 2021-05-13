package main

import (
	"flag"
	"fmt"
	"log"

	_ "github.com/powerpuffpenguin/sessionid_server/assets/document/statik"
	"github.com/powerpuffpenguin/sessionid_server/configure"
	_ "github.com/powerpuffpenguin/sessionid_server/gmodule/manager"
	_ "github.com/powerpuffpenguin/sessionid_server/gmodule/provider"
	"github.com/powerpuffpenguin/sessionid_server/logger"
	"github.com/powerpuffpenguin/sessionid_server/server"
	"github.com/powerpuffpenguin/sessionid_server/system"
	"github.com/powerpuffpenguin/sessionid_server/version"
	"go.uber.org/zap"
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
		e = logger.Init(&configure.DefaultConfigure().Logger)
		if e != nil {
			log.Fatalln(e)
		}
		system.Init()

		serverCnf := configure.DefaultConfigure().Server
		srv, e := server.NewServer(serverCnf.Addr)
		if e != nil {
			log.Fatalln(e)
		}
		if serverCnf.CertFile == `` && serverCnf.KeyFile == `` {
			logger.Logger.Info(`h2c work`,
				zap.String(`addr`, serverCnf.Addr),
			)
			e = srv.Serve()
		} else {
			logger.Logger.Info(`h2 work`,
				zap.String(`addr`, serverCnf.Addr),
			)
			e = srv.ServeTLS(serverCnf.CertFile, serverCnf.KeyFile)
		}
		if e != nil {
			log.Fatalln(e)
		}
	}
}
