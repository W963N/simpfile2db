package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/naoina/toml"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type dbConfig struct {
	Root   string
	Output string
	Dbpath string
	Env    map[string]EnvInfo
}

type EnvInfo struct {
	Name        string
	Description string
	Ext         string
}

func (cfg *dbConfig) loadConf(file io.Reader) error {
	return toml.NewDecoder(file).Decode(&cfg)
}

var reg = `{{ red "Register:" }} {{ bar . "W9" "6" (cycle . "6" "6" "6" "6" ) "." "3N"}} {{speed . | rndcolor }} {{percent .}} {{string . "gst" | green}} {{string . "bst" | blue}}`
var out = `{{ red "Output:" }} {{ bar . "W9" "6" (cycle . "6" "6" "6" "6" ) "." "3N"}} {{speed . | rndcolor }} {{percent .}} {{string . "gst" | green}} {{string . "bst" | blue}}`
var sch = `{{ red "Search:" }} {{ bar . "W9" "6" (cycle . "6" "6" "6" "6" ) "." "3N"}} {{speed . | rndcolor }} {{percent .}} {{string . "gst" | green}} {{string . "bst" | blue}}`

func (cfg *dbConfig) createDBPath(tomlabel string) (string, error) {
	_, ok := cfg.Env[tomlabel]
	if !ok {
		return "", errors.New("DB doesn't exist.")
	}
	return cfg.Dbpath + "/" +
		cfg.Env[tomlabel].Name, nil
}

func (cfg *dbConfig) getExt(tomlabel string) (string, error) {
	_, ok := cfg.Env[tomlabel]
	if !ok {
		return "", errors.New("Don't get Ext conf")
	}
	return cfg.Env[tomlabel].Ext, nil
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

func (cfg *dbConfig) regDB(db *leveldb.DB) error {
	count := int64(5)
	bar := pb.ProgressBarTemplate(reg).Start64(count)
	bar.Set("bst", "fin").Set("gst", TITLE).
		SetMaxWidth(80).SetRefreshRate(time.Second)
	if err := bar.Err(); err != nil {
		return err
	}

	ext, err := cfg.getExt(target_flag)
	if err != nil {
		return err
	}
	bar.Increment()

	if err := chkExt(file_flag, ext); err != nil {
		return err
	}
	bar.Increment()

	reg_file := cfg.Root + "/" + file_flag
	bar.Increment()

	bytes, err := ioutil.ReadFile(reg_file)
	if err != nil {
		return err
	}
	bar.Increment()

	key := getFileNameWithoutExt(reg_file)
	if err := db.Put([]byte(key), bytes, nil); err != nil {
		return err
	}
	bar.Increment()
	bar.Finish()

	return err
}

func (cfg *dbConfig) outDB(db *leveldb.DB) error {
	count := int64(3)
	bar := pb.ProgressBarTemplate(out).Start64(count)
	bar.Set("bst", "fin").Set("gst", TITLE).
		SetMaxWidth(80).SetRefreshRate(time.Second)
	if err := bar.Err(); err != nil {
		return err
	}

	data, err := db.Get([]byte(key_flag), nil)
	if err != nil {
		return err
	}
	bar.Increment()

	out := cfg.Output + "/" + key_flag + cfg.Env[target_flag].Ext
	wf, err := os.Create(out)
	if err != nil {
		return err
	}
	defer wf.Close()
	bar.Increment()

	_, err = wf.Write(data)
	if err != nil {
		return err
	}
	bar.Increment()
	bar.Finish()

	return err
}

func (cfg *dbConfig) searchDB(db *leveldb.DB) error {
	var keys []string
	iter := db.NewIterator(util.BytesPrefix([]byte(key_flag)), nil)
	for iter.Next() {
		key := iter.Key()
		//value := iter.Value()
		keys = append(keys, string(key))
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return err
	}

	count := int64(len(keys))
	bar := pb.ProgressBarTemplate(sch).Start64(count)
	bar.Set("bst", "fin").Set("gst", TITLE).
		SetMaxWidth(80).SetRefreshRate(time.Second)
	if err := bar.Err(); err != nil {
		return err
	}

	for _, key := range keys {
		bar.Increment()
		fmt.Println(key)
	}
	bar.Finish()

	return nil
}

func (cfg *dbConfig) deleteDB(db *leveldb.DB) error {
	if err := db.Delete([]byte(key_flag), nil); err != nil {
		return err
	}

	return nil
}

func (cfg *dbConfig) execDB() error {
	dbname, err := cfg.createDBPath(target_flag)
	if err != nil {
		return err
	}

	db, err := leveldb.OpenFile(dbname, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	if reg_flag {
		if err := cfg.regDB(db); err != nil {
			return err
		}
	} else if search_flag {
		if err := cfg.searchDB(db); err != nil {
			return err
		}
	} else if del_flag {
		if err := cfg.deleteDB(db); err != nil {
			return err
		}
	} else if output_flag {
		if err := cfg.outDB(db); err != nil {
			return err
		}
	} else {
		panic("Unknown flag")
	}

	return err
}
