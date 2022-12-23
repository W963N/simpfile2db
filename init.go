package main

import (
	"errors"
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func init() {
	flag.CommandLine.Init("DB REG T00l", flag.ContinueOnError)
	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "\nW963N Memory ~~ %s\n", flag.CommandLine.Name())
		fmt.Fprintf(o, "\n@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n")
		fmt.Fprintf(o, "@@@@@@@@@@@@@  .@@@@@@@@@.  @@@@@@@@@@@@\n")
		fmt.Fprintf(o, "@@@@@@@@@ (@@@@@@@@@@@@@@@@@@@) @@@@@@@@\n")
		fmt.Fprintf(o, "@@@@@@ @@ @@@@             @@@@ @@ @@@@@\n")
		fmt.Fprintf(o, "@@@@ @@@@@ @ @@@         @@@ @ @@@@ @@@@\n")
		fmt.Fprintf(o, "@@@ @@@@@@ @@@@@@       @@@@@@ @@@@@ @@@\n")
		fmt.Fprintf(o, "@@ @@@@@@,@@@@@@@       @@@@@@@,@@@@@ @@\n")
		fmt.Fprintf(o, "@@ @@@@@@ @@@@@@@@     @@@@@@@@ @@@@@@ @\n")
		fmt.Fprintf(o, "@       @@@      @     @      @@@      @\n")
		fmt.Fprintf(o, "@@ @@@ @@@@@@   @@@   @@@   @@@@@@ @@ @@\n")
		fmt.Fprintf(o, "@@@ @@ @@@@@@@@@@@@   @@@@@@@@@@@@ @@ @@\n")
		fmt.Fprintf(o, "@@@@  @@@@@@@@@@@@@@ @@@@@@@@@@@@@ @ @@@\n")
		fmt.Fprintf(o, "@@@@@  @@@@@@@@@@@@@ @@@@@@@@@@@@   @@@@\n")
		fmt.Fprintf(o, "@@@@@@@   @@@@@@@  @@@  @@@@@@@   @@@@@@\n")
		fmt.Fprintf(o, "@@@@@@@@@@ @@@@@@@@@@@@@@@@@@@ @@@@@@@@@\n")
		fmt.Fprintf(o, "@@@@@@@@@@@@@@@     @     @@@@@@@@@@@@@@\n\n")
		fmt.Fprintf(o, "\nUsage: \n")
		fmt.Fprintf(o, "  simpfile2db -t [DB NAME] -r -f [KEY(FILE NAME)].[Ext]\n")
		fmt.Fprintf(o, "  simpfile2db -t [DB NAME] -s -k [KEY]\n")
		fmt.Fprintf(o, "  simpfile2db -t [DB NAME] -d -k [KEY]\n")
		fmt.Fprintf(o, "  simpfile2db -t [DB NAME] -o -k [KEY]\n")
		fmt.Fprintf(o, "\nOptions: \n")
		flag.PrintDefaults()
		fmt.Fprintf(o, "\nHobbyright 2022 walnut 🐿🐿🐿 .\n\n")
	}
	flag.StringVar(&env_flag, "e", "./env.toml", "path of env.toml.")
	flag.StringVar(&target_flag, "t", "", "Use db name.")
	flag.BoolVar(&reg_flag, "r", false, "Register file.")
	flag.BoolVar(&search_flag, "s", false, "Search file.")
	flag.BoolVar(&del_flag, "d", false, "Delete file.")
	flag.BoolVar(&output_flag, "o", false, "Output file.")
	flag.StringVar(&file_flag, "f", "", "file path")
	flag.StringVar(&key_flag, "k", "", "key")
	flag.StringVar(&verbose_flag, "v", "error", "Select types(info, warn, error).")
}

var (
	env_flag     string
	reg_flag     bool
	search_flag  bool
	del_flag     bool
	output_flag  bool
	target_flag  string
	file_flag    string
	key_flag     string
	verbose_flag string
)

func chkFlag(target_flag string, output_flag, reg_flag, search_flag, del_flag bool) bool {
	if target_flag == "" {
		log.Error("Please enter db name.")
		return true
	}
	if !(reg_flag || search_flag || del_flag || output_flag) {
		log.Error("Please select one of the options(-r -s -d -o).")
		return true
	}
	if reg_flag && search_flag && del_flag && output_flag {
		log.Error("Please select one of the options(-r -s -d -o).")
		return true
	}
	if reg_flag && search_flag && del_flag && !output_flag {
		log.Error("Please select one of the options(-r -s -d).")
		return true
	}
	if reg_flag && search_flag && !del_flag && output_flag {
		log.Error("Please select one of the options(-r -s -o).")
		return true
	}
	if reg_flag && !search_flag && del_flag && output_flag {
		log.Error("Please select one of the options(-r -d -o).")
		return true
	}
	if !reg_flag && search_flag && del_flag && output_flag {
		log.Error("Please select one of the options(-s -d -o).")
		return true
	}
	if reg_flag && search_flag {
		log.Error("Please select either option(-r or -s).")
		return true
	}
	if reg_flag && del_flag {
		log.Error("Please select either option(-r or -d).")
		return true
	}
	if reg_flag && output_flag {
		log.Error("Please select either option(-r or -o).")
		return true
	}
	if search_flag && del_flag {
		log.Error("Please select either option(-s or -d).")
		return true
	}
	if search_flag && output_flag {
		log.Error("Please select either option(-s or -o).")
		return true
	}
	if del_flag && output_flag {
		log.Error("Please select either option(-d or -o).")
		return true
	}
	if reg_flag && file_flag == "" {
		log.Error("Please select registration file path.")
		return true
	}
	if search_flag && key_flag == "" {
		log.Error("Please enter search key.")
		return true
	}
	if del_flag && key_flag == "" {
		log.Error("Please enter delete key.")
		return true
	}
	if output_flag && key_flag == "" {
		log.Error("Please enter output key.")
		return true
	}
	return false
}

func changeLogLevel(level string) error {
	switch level {
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		return errors.New("Don't match level.(info, warn, error)")
	}
	return nil
}
