package config

import (
	"flag"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
)

type config struct {
	Host          string
	Port          int
	HttpsEnabled  bool
	CertPath      string
	KeyPath       string
	DatabaseUrl   string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

var Conf config

func init() {
	var confPath string
	flag.StringVar(&confPath, "conf", "", "Config File path")
	flag.Parse()
	if len(confPath) != 0 {
		var bytes []byte
		var err error
		if bytes, err = ioutil.ReadFile(confPath); err != nil {
			log.Fatal(err)
		}
		if _, err = toml.Decode(string(bytes), &Conf); err != nil {
			log.Fatal(err)
		}
	}
}
