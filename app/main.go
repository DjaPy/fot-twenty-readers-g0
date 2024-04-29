package main

import (
	"context"
	"fmt"
	"github.com/DjaPy/fot-twenty-readers-go/app/api"
	"github.com/DjaPy/fot-twenty-readers-go/app/config"
	log "github.com/go-pkgz/lgr"
	"github.com/jessevdk/go-flags"
	"os"
)

type options struct {
	Port int    `short:"p" long:"port" description:"port to listen" default:"8080"`
	Conf string `short:"f" long:"conf" env:"FM_CONF" default:"for_twenty_readers.yml" description:"config file (yml)"`
	Dbg  bool   `long:"dbg" env:"DEBUG" description:"debug mode"`
}

var revision = "local"

func main() {
	fmt.Printf("For twenty readers %s\n", revision)
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(1)
	}
	setupLog(opts.Dbg)

	var conf = &config.Conf{}

	server := api.Server{
		Version: revision,
		Conf:    *conf,
	}
	server.Run(context.Background(), opts.Port)
}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}
