package main

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/google/goterm/term"
	log "github.com/sirupsen/logrus"
)

const (
	TITLE = "DB REG T00l"
)

func clyFprintf(w io.Writer, format string) {
	fmt.Fprintf(w, term.Yellowf(format))
}

func clrFprintf(w io.Writer, format string) {
	fmt.Fprintf(w, term.Redf(format))
}

func clbFprintf(w io.Writer, format string) {
	fmt.Fprintf(w, term.Bluef(format))
}

func init() {
	flag.CommandLine.Init(TITLE, flag.ContinueOnError)
	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "\nW963N Memory ~~ %s\n", flag.CommandLine.Name())
		clyFprintf(o, "\n@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n")
		clyFprintf(o, "@@@@@@@@@@@@@  .@@@@@@@@@.  @@@@@@@@@@@@\n")
		clyFprintf(o, "@@@@@@@@@ (@@@@@@@@@@@@@@@@@@@) @@@@@@@@\n")
		clyFprintf(o, "@@@@@@ @@ @@@@             @@@@ @@ @@@@@\n")
		clyFprintf(o, "@@@@ @@@@@ @ @@@         @@@ @ @@@@ @@@@\n")
		clyFprintf(o, "@@@ @@@@@@ @@@@@@       @@@@@@ @@@@@ @@@\n")
		clyFprintf(o, "@@ @@@@@@,@@@@@@@       @@@@@@@,@@@@@ @@\n")
		clyFprintf(o, "@@ @@@@@@ @@@@@@@@     @@@@@@@@ @@@@@@ @\n")
		clyFprintf(o, "@       @@@      @     @      @@@      @\n")
		clyFprintf(o, "@@ @@@ @@@@@@   @@@   @@@   @@@@@@ @@ @@\n")
		clyFprintf(o, "@@@ @@ @@@@@@@@@@@@   @@@@@@@@@@@@ @@ @@\n")
		clyFprintf(o, "@@@@  @@@@@@@@@@@@@@ @@@@@@@@@@@@@ @ @@@\n")
		clyFprintf(o, "@@@@@  @@@@@@@@@@@@@ @@@@@@@@@@@@   @@@@\n")
		clyFprintf(o, "@@@@@@@   @@@@@@@  @@@  @@@@@@@   @@@@@@\n")
		clyFprintf(o, "@@@@@@@@@@ @@@@@@@@@@@@@@@@@@@ @@@@@@@@@\n")
		clyFprintf(o, "@@@@@@@@@@@@@@@     @     @@@@@@@@@@@@@@\n\n")
		clrFprintf(o, "\nUsage: \n")
		fmt.Fprintf(o, "  simpfile2db -t [DB NAME] -r -f [KEY(FILE NAME)].[Ext]\n")
		fmt.Fprintf(o, "  simpfile2db -t [DB NAME] -s -k [KEY]\n")
		fmt.Fprintf(o, "  simpfile2db -t [DB NAME] -d -k [KEY]\n")
		fmt.Fprintf(o, "  simpfile2db -t [DB NAME] -o -k [KEY]\n")
		clrFprintf(o, "\nOptions: \n")
		flag.PrintDefaults()
		clbFprintf(o, "\nHobbyright 2022 walnut 🐿🐿🐿 .\n\n")
	}
	flag.StringVar(&env_flag, term.Greenf("e"), "./env.toml", "path of env.toml.")
	flag.StringVar(&target_flag, term.Greenf("t"), "", "Use db name.")
	flag.BoolVar(&reg_flag, term.Greenf("r"), false, "Register file.")
	flag.BoolVar(&search_flag, term.Greenf("s"), false, "Search file.")
	flag.BoolVar(&del_flag, term.Greenf("d"), false, "Delete file.")
	flag.BoolVar(&output_flag, term.Greenf("o"), false, "Output file.")
	flag.StringVar(&file_flag, term.Greenf("f"), "", "file path")
	flag.StringVar(&key_flag, term.Greenf("k"), "", "key")
	flag.StringVar(&verbose_flag, term.Greenf("v"), "error", "Select types(info, warn, error).")
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
