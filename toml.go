package main

import (
	"io"

	"github.com/naoina/toml"
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

func loadConf(file io.Reader, config *dbConfig) error {
	return toml.NewDecoder(file).Decode(&config)
}
