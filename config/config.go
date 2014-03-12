package config

import (
	"encoding/json"
	"os"
	"sort"
)

var Config = &struct {
	Domains Domains
}{
	Domains: make(Domains, 0, 64),
}

func Add(d *Domain) {
	Config.Domains = append(Config.Domains, d)
	sort.Sort(Config.Domains)
}

func Get(alias string) (*Domain, error) {
	return Config.Domains.Get(alias)
}

func Load(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	d := json.NewDecoder(f)
	return d.Decode(Config)
}

func Save(filename string) (err error) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	e := json.NewEncoder(f)
	return e.Encode(Config)
}
