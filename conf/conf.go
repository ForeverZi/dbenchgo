package conf

import (
	"flag"
	"github.com/spf13/viper"
	"log"
)

var confFile string

func init() {
	flag.StringVar(&confFile, "f", "conf.yaml", "-f conf.yaml")
	flag.Parse()
	viper.SetConfigFile(confFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	log.Println("conf loaded...")
}
