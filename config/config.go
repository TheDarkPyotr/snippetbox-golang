package config

import (
	"log"
	"time"

	"github.com/joeshaw/envdecode"
)

type dbConf struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DbName   string `env:"DB_NAME,required"`
}

type TLSServerInfo struct {
	CertificatePath string `env:"TLS_CERT_PATH,required"`
	PrivateKeyPath  string `env:"TLS_PKEY_PATH,required"`
}

type CookieInfo struct {
	Secret32 string `env:"COOKIE_PKEY,required"`
}

type HTTPDefaultParam struct {
	Addr         string        `env:"HTTP_SERVER_PORT,required"`
	IdleTimeout  time.Duration `env:"HTTP_SERVER_TIMEOUT_IDLE,required"`
	ReadTimeout  time.Duration `env:"HTTP_SERVER_TIMEOUT_READ,required"`
	WriteTimeout time.Duration `env:"HTTP_SERVER_TIMEOUT_WRITE,required"`
}

type Specification struct {
	ServerParams HTTPDefaultParam

	StaticDir     string `env:"STATIC_PATH,required"`
	DbConf        dbConf
	TLS           TLSServerInfo
	CookieSetting CookieInfo
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
