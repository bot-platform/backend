package main

type Config struct {
	Port   int
	DbHost string
	DbName string
}

func NewConfig() *Config {
	return &Config{
		Port:   1323,
		DbHost: "http://localhost:27017",
		DbName: "bot-platform",
	}
}
