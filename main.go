package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func die(args ...interface{}) {
	log.Error(args)
	os.Exit(1)
}

func main() {
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		os.Exit(0)
	}
	if chkFlag(target_flag, output_flag, reg_flag, search_flag, del_flag) {
		die("Check Flag error")
	}
	if err := changeLogLevel(verbose_flag); err != nil {
		die(err)
	}

	file, err := os.Open(env_flag)
	if err != nil {
		die(err)
	}
	defer file.Close()

	cfg := &dbConfig{}
	if err := cfg.loadConf(file); err != nil {
		die(err)
	}

	log.Info("[ASM ENV]:name: ", cfg.Env[target_flag].Name)
	log.Info("[ASM ENV]:root path: ", cfg.Root)
	log.Info("[ASM ENV]:dbpath: ", cfg.Dbpath)

	if err := cfg.execDB(); err != nil {
		die(err)
	}
}
