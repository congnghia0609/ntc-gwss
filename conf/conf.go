package conf

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper (external lib) and
// returns the configuration struct.
func Init(env string) {
	log.Printf("============== Config Init Environment: %s ==============", env)
	var err error
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName(env)
	v.AddConfigPath("../conf/")
	v.AddConfigPath("conf/")
	err = v.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
	config = v
}

func relativePath(basedir string, path *string) {
	p := *path
	if p != "" && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

func GetConfig() *viper.Viper {
	return config
}
