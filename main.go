package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var configFile = flag.String("conf", "Cakefile.yaml", "cakeserver -conf=<file>")

func main() {
	flag.Parse()
	conf, err := getConfig(*configFile)
	if err != nil {
		log.Fatalln(err)
	}

	if conf.Addr == "" {
		log.Fatalln("addr not presented.")
	}
	if conf.Secret == "" {
		log.Fatalln("secret not presented.")
	}

	Commands = conf.Commands
	if Commands == nil {
		Commands = []Command{}
	}
	serve(conf.Addr, conf.Secret)
}

type Config struct {
	Addr     string    `yaml:"addr"`
	Secret   string    `yaml:"secret"`
	Commands []Command `yaml:"commands"`
}

func getConfig(file string) (*Config, error) {
	confFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	yml, err := ioutil.ReadAll(confFile)
	if err != nil {
		return nil, err
	}

	configs := Config{}
	if err := yaml.Unmarshal(yml, &configs); err != nil {
		return nil, err
	}

	return &configs, nil
}
