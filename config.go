package main

import (
	"gopkg.in/ini.v1"
	"runtime"
	"time"
)

type Config struct {
	httpConfig   *HttpConfig
	pingerConfig *PingerConfig
}

type HttpConfig struct {
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	idleTimeout  time.Duration
}

type PingerConfig struct {
	numbersOfWorker int
	httpPort        string
}

func createConfig() *Config {
	return &Config{&HttpConfig{}, &PingerConfig{}}
}

func (config *Config) parse() {
	filepath := "config/config.ini"
	cfg, err := ini.Load(filepath)
	if err != nil {
		panic(err)
	}

	port := cfg.Section("HTTP_SERVER").Key("Port").String()
	if port == "" {
		port = "8080"
	}
	config.httpConfig.port = port

	readTimeout, err := cfg.Section("HTTP_SERVER").Key("ReadTimeout").Int()
	if err != nil {
		readTimeout = 60
	}
	config.httpConfig.readTimeout = time.Duration(readTimeout)

	writeTimeout, err := cfg.Section("HTTP_SERVER").Key("WriteTimeout").Int()
	if err != nil {
		writeTimeout = 60
	}
	config.httpConfig.writeTimeout = time.Duration(writeTimeout)

	idleTimeout, err := cfg.Section("HTTP_SERVER").Key("IdleTimeout").Int()
	if err != nil {
		idleTimeout = 60
	}
	config.httpConfig.idleTimeout = time.Duration(idleTimeout)

	numbersOfWorker, err := cfg.Section("PINGER").Key("NumbersOfWorker").Int()
	if err != nil {
		numbersOfWorker = runtime.NumCPU() * 2
	}
	config.pingerConfig.numbersOfWorker = numbersOfWorker

	responseTimeout, err := cfg.Section("PINGER").Key("MaxRTT").Int()
	if err != nil {
		responseTimeout = 50
	}
	ResponseTimeout = responseTimeout
}

func (config *Config) erase() {
	config.httpConfig = nil
	config.pingerConfig = nil
}
