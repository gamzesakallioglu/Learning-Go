package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	ServerConfig ServerConfig
	JWTConfig    JWTConfig
	DBConfig     DBConfig
	Logger       Logger
}

type ServerConfig struct {
	AppVersion       string
	Mode             string
	RoutePrefix      string
	Debug            bool
	Port             string
	TimeoutSecs      int64
	ReadTimeoutSecs  int64
	WriteTimeoutSecs int64
}

type JWTConfig struct {
	SessionTime int64
	SecretKey   string
}

type DBConfig struct {
	DataSourceName string
	MaxOpen        int
	MaxIdle        int
	MaxLifetime    int
}

type Logger struct {
	Development bool
	Encoding    string
	Level       string
}

func LoadConfig(filename string) (*Config, error) {
	vpr := viper.New()

	vpr.SetConfigName(filename) // set config file's name
	vpr.AddConfigPath(".")      // look for config in the working directory
	vpr.AutomaticEnv()

	if err := vpr.ReadInConfig(); err != nil {

		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file could not be found")
		}
		return nil, err

	}

	var c Config
	if err := vpr.Unmarshal(&c); err != nil {
		log.Printf("could not decode into struct, %v", err)
		return nil, err
	}

	return &c, nil

}
