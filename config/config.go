package config

import (
	"log"

	"github.com/joeshaw/envdecode"
)

type dbConf struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DbName   string `env:"DB_NAME,required"`
}

type Specification struct {
	Addr      string `env:"HTTP_SERVER_PORT,required"`
	StaticDir string `env:"STATIC_PATH,required"`
	DbConf    dbConf
}

func ServerConfiguration() *Specification {

	var conf Specification
	if err := envdecode.StrictDecode(&conf); err != nil {
		log.Fatalf("Error decoding server conf: %s", err)

	}

	return &conf

}

func AppConfig() *Specification {
	var c Specification
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}

func DbConfig() *dbConf {
	var c dbConf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	} else {
		log.Printf("DB Name %s", c.DbName)
	}

	return &c
}
