package utils

import (
	"errors"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlStruct struct {
	Server   Server   `yaml:"server"`
	Account  Account  `yaml:"account"`
	DataUser DataUser `yaml:"data_user"`
	Config   Config   `yaml:"config"`
}

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Account struct {
	Alias          string `yaml:"alias"`
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	UsernameSender string `yaml:"username_sender"`
}

type DataUser struct {
	EmailAdmin string `yaml:"email_admin"`
	Homepage   string `yaml:"homepage"`
	Phone      string `yaml:"phone"`
	Guide      Guide  `yaml:"guide"`
}

type Guide struct {
	AdminMerchant string `yaml:"admin_merchant"`
	Merchant      string `yaml:"merchant"`
}

type Config struct {
	Port    string `yaml:"port"`
	Token   string `yaml:"token"`
	Storage string `yaml:"storage"`
}

func OpenYAML() (*YamlStruct, error) {
	var pathFile string
	if _, err := os.Stat("config.yaml"); err == nil {
		pathFile = "config.yaml"
	} else if _, err := os.Stat("config.yml"); err == nil {
		pathFile = "config.yml"
	} else {
		return nil, errors.New("config yaml file not found")
	}

	data, err := os.ReadFile(pathFile)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var config YamlStruct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &config, nil
}
