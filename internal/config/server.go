package config

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Port           string       `ini:"port"`
	LogLevel       logrus.Level `ini:"log_level"`
	HistoryStorage string       `ini:"history_storage"`
}

type DataBaseConfig struct {
	Host     string `ini:"host"`
	Port     string `ini:"port"`
	User     string `ini:"user"`
	Password string `ini:"password"`
	DBName   string `ini:"db_name"`
}

type Config struct {
	ServerConfig   ServerConfig   `ini:"server"`
	DataBaseConfig DataBaseConfig `ini:"database"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := ini.Load(path)
	if err != nil {
		return nil, err
	}
	conf := &Config{}
	if err := f.MapTo(conf); err != nil {
		return nil, err
	}
	return conf, nil
}
