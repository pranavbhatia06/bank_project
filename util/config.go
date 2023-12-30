package util

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Application App      `json:"application"`
	Database    Database `json:"database"`
}
type App struct {
	AppEnv    string
	LISTEN_IP string
}
type Database struct {
	Dialect        string
	Port           int
	Host           string
	DbName         string
	SslMode        string
	DataSourceName string
}

func LoadConfig(path string) (Config, error) {
	fileName := "default"
	appEnv := os.Getenv("APP_ENV")

	if appEnv != "" {
		fileName = appEnv
	}
	viper.SetConfigName(fileName)
	viper.SetConfigType("toml")

	viper.AddConfigPath(path)
	viper.AutomaticEnv()

	cnf := Config{}

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading configuration file:", err)
		return Config{}, err
	}
	if err := viper.Unmarshal(&cnf); err != nil {
		fmt.Println("Error unmarshalling config:", err)
		return Config{}, err
	}
	cnf.Database.DataSourceName = fmt.Sprintf("postgresql://root:password@%s:%d/%s?sslmode=%s", cnf.Database.Host, cnf.Database.Port, cnf.Database.DbName, cnf.Database.SslMode)

	return cnf, nil
}
