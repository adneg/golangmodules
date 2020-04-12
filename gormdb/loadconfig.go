package gormdb

import (
	"encoding/json"
	"os"

	"github.com/adneg/golangmodules/logtrace"
)

func loadconfig(place string) (conf Configuration) {

	file, err := os.Open(place)
	if err != nil {
		logtrace.Error.Fatalln(err.Error())

	}
	defer file.Close()
	decoder := json.NewDecoder(file)

	conf = Configuration{}
	err = decoder.Decode(&conf)

	if err != nil {
		logtrace.Error.Fatal(err.Error())

	}
	logtrace.Info.Println("DB CONFIG LOADED")
	confstr, _ := json.Marshal(conf)
	logtrace.Info.Println(string(confstr))
	return conf

}
