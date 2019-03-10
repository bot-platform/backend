package main

type Config struct {
	Port      int
	DbHost    string
	DbName    string
	JWTSecret interface{}
}

func NewConfig() *Config {
	return &Config{
		Port:      1323,
		DbHost:    "http://localhost:27017",
		DbName:    "bot-platform",
		JWTSecret: []byte("secret"),
	}
}
