package configuration

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

/*
	- Load yml file  
    - Read it into struct
	- store it into map with host as key and servers as value array

	[{}, {}, {}]
*/

type Config struct {
	Hosts map[string][]string `yaml:"Hosts"`
}

func read(filename string) []byte{
	data,err := ioutil.ReadFile(filename)
	if err != nil{
		log.Fatal(err)
	}
	return data
}

func Load(filename string) Config{
	var config Config

	stream := read(filename)
	yaml.Unmarshal(stream, &config)

	return Config
}