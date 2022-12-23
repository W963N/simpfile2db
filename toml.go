package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/naoina/toml"
	log "github.com/sirupsen/logrus"
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
	ext, err := cfg.getExt(target_flag)
	if err != nil {
		return err
	}
	if err := chkExt(file_flag, ext); err != nil {
		return err
	}
	reg_file := cfg.Root + "/" + file_flag
	log.Info("[REG FILE]:path: ", reg_file)

	bytes, err := ioutil.ReadFile(reg_file)
	if err != nil {
		return err
	}
	key := getFileNameWithoutExt(reg_file)
	log.Info("[REG FILE]:key: ", key)
	if err := db.Put([]byte(key), bytes, nil); err != nil {
		return err
	}
	log.Info("[REG FILE]:success!!")

	return err
}

func (cfg *dbConfig) outDB(db *leveldb.DB) error {
	data, err := db.Get([]byte(key_flag), nil)
	if err != nil {
		return err
	}
	log.Info("byte: ", data)

	out := cfg.Output + "/" + key_flag + cfg.Env[target_flag].Ext
	log.Info("[OUTPUT]:output path: ", out)
	wf, err := os.Create(out)
	if err != nil {
		return err
	}
	defer wf.Close()

	_, err = wf.Write(data)
	if err != nil {
		return err
	}
	log.Info("[OUTPUT]: success!!")

	return err
}

func (cfg *dbConfig) searchDB(db *leveldb.DB) error {
	iter := db.NewIterator(util.BytesPrefix([]byte(key_flag)), nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println(string(key))
		log.Info("[SEARCH]:key: ", string(key))
		log.Info("[SEARCH]:value: ", hex.EncodeToString(value))
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		return err
	}

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
	log.Info("[TARGT DB]:path: ", dbname)

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
