package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
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

func createDBPath(tomlabel string, config dbConfig) (string, error) {
	_, ok := config.Env[tomlabel]
	if !ok {
		return "", errors.New("DB doesn't exist.")
	}
	return config.Dbpath + "/" +
		config.Env[tomlabel].Name, nil
}

func getExt(tomlabel string, config dbConfig) (string, error) {
	_, ok := config.Env[tomlabel]
	if !ok {
		return "", errors.New("DB doesn't exist.")
	}
	return config.Env[tomlabel].Ext, nil
}

func getFileNameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}

func chkExt(file string, ext string) error {
	if filepath.Ext(file) != ext {
		return errors.New("Don't match Ext")
	}
	return nil
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

func main() {
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		os.Exit(0)
	}
	if chkFlag(target_flag, output_flag, reg_flag, search_flag, del_flag) {
		os.Exit(2)
	}
	if err := changeLogLevel(verbose_flag); err != nil {
		log.Error(err)
		os.Exit(2)
	}

	file, err := os.Open(env_flag)
	if err != nil {
		log.Error("Don't open env file")
		os.Exit(1)
	}
	defer file.Close()
	var config dbConfig
	if err := loadConf(file, &config); err != nil {
		log.Error("Failed to decode toml file")
		os.Exit(1)
	}

	log.Info("[ASM ENV]:name: ", config.Env[target_flag].Name)
	log.Info("[ASM ENV]:root path: ", config.Root)
	log.Info("[ASM ENV]:dbpath: ", config.Dbpath)

	dbname, err := createDBPath(target_flag, config)
	if err != nil {
		log.Error(err)
		os.Exit(3)
	}
	log.Info("[TARGT DB]:path: ", dbname)

	db, err := leveldb.OpenFile(dbname, nil)
	if err != nil {
		log.Error(err)
		os.Exit(4)
	}
	defer db.Close()

	if reg_flag {
		ext, err := getExt(target_flag, config)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		err = chkExt(file_flag, ext)
		if err != nil {
			log.Error(err)
			os.Exit(5)
		}
		reg_file := config.Root + "/" + file_flag
		log.Info("[REG FILE]:path: ", reg_file)

		bytes, err := ioutil.ReadFile(reg_file)
		if err != nil {
			log.Error("[REG FILE]:Don't open file")
			os.Exit(5)
		}
		key := getFileNameWithoutExt(reg_file)
		log.Info("[REG FILE]:key: ", key)
		err = db.Put([]byte(key), bytes, nil)
		if err != nil {
			log.Error("[REG FILE]: Registration failed.")
			os.Exit(5)
		}
		log.Info("[REG FILE]:success!!")
	} else if search_flag {
		iter := db.NewIterator(util.BytesPrefix([]byte(key_flag)), nil)
		for iter.Next() {
			key := iter.Key()
			value := iter.Value()
			fmt.Println(string(key))
			log.Info("[SEARCH]:key: ", string(key))
			log.Info("[SEARCH]:value: ", hex.EncodeToString(value))
		}
		iter.Release()
		err = iter.Error()
		if err != nil {
			log.Error("[SEARCH]: Failed iteration.")
			os.Exit(6)
		}
	} else if del_flag {
		err = db.Delete([]byte(key_flag), nil)
		if err != nil {
			log.Error("[DELETE]: Failed.")
			os.Exit(7)
		}
	} else if output_flag {
		data, err := db.Get([]byte(key_flag), nil)
		if err != nil {
			log.Error("[OUTPUT]: ", err)
			os.Exit(8)
		}
		log.Info("byte: ", data)

		out := config.Output + "/" + key_flag + config.Env[target_flag].Ext
		log.Info("[OUTPUT]:output path: ", out)
		wf, err := os.Create(out)
		if err != nil {
			log.Error(err)
			os.Exit(8)
		}
		defer wf.Close()
		_, err = wf.Write(data)
		if err != nil {
			log.Error(err)
			os.Exit(8)
		}
		log.Info("[OUTPUT]: success!!")
	} else {
		log.Error("Unknown flag")
		os.Exit(2)
	}
}
