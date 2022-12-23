package main

import (
	"flag"
	"fmt"
	"os"
)

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args)
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

	file, err := os.Open(env_flag)
	if err != nil {
		die("Open Env file error %v:", err)
	}
	defer file.Close()

	cfg := &dbConfig{}
	if err := cfg.loadConf(file); err != nil {
		die("Struct Config error %v:", err)
	}

	if err := cfg.execDB(); err != nil {
		die("DB operation error %v:", err)
	}
}
