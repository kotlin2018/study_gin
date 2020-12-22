package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// 配置结构体
type Server struct {

}


func Config() Server {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		panic(err)
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	s := Server{}
	err = yaml.Unmarshal(bytes, &s)
	if err != nil {
		panic(err)
	}
	return s
}
