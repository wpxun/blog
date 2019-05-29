package service

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type confstruct struct {
	Redis struct {
		Address     string 			`yaml:"Address"`
		Database    int    			`yaml:"Database"`
		DialTimeout time.Duration   `yaml:"DialTimeout"`
		Network     string 			`yaml:"Network"`
		Password    string 			`yaml:"Password"`
	} `yaml:"Redis"`

	MathRPC struct {
		Address     string 			`yaml:"Address"`
	} `yaml:"MathRPC"`
}

var Conf confstruct

func init() {
	GetYaml("config", &Conf)
}

func GetYaml(filename string, out interface{}) {
	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("%s.yaml", filename))
	if err != nil {
		fmt.Println("Read config file error:", err.Error())
	}

	err = yaml.Unmarshal(yamlFile, out)

	if (err != nil) {
		fmt.Println("Unmarshal config file error:", err.Error())
	}
}
