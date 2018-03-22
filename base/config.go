package base

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os/user"
	"path"
)

var config Config

// Mail config for mail(smtp)
type Mail struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

// GetConfig get config
func GetConfig() Config {
	return config
}

// Config : config schema
type Config struct {
	ProjectName string           `json:"name"`
	Mail        Mail             `json:"mail"`
	Plans       *json.RawMessage `json:"plans"`
}

// LoadConfig according to path
func LoadConfig(pa string) {
	if pa == "" {
		usr, _ := user.Current()
		pa = path.Join(usr.HomeDir, ".hc", "config.json")
	}
	conf, err := ioutil.ReadFile(pa)
	if err != nil {
		log.Panic("Failed to read config file.")
	}
	err = json.Unmarshal(conf, &config)
	if err != nil {
		log.Panic("Config file is broken.")
	}
}

var pa = flag.String("config", "", "path of config.json")

func init() {
	flag.Parse()
	LoadConfig(*pa)
}
