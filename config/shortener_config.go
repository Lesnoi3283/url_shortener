package config

import (
	"flag"
	"os"
)

const DefaultBaseAddress = "http://localhost:8080"
const DefaultServerAddress = "localhost:8080"
const DefaultLogLevel = "info"
const DefaultFileStoragePath = "/tmp/short-url-db.json"
const DefaultDBConnectionString = ""

// Config is a struct with configuration params.
type Config struct {
	BaseAddress     string
	ServerAddress   string
	LogLevel        string
	FileStoragePath string
	DBConnString    string
}

// Configure reads configuration params from command line args, environmental variables and DefaultConstParams.
// And writes them into a Config struct.
func (c *Config) Configure() {
	flag.StringVar(&(c.ServerAddress), "a", DefaultServerAddress, "Address where server will work. Example: \"localhost:8080\".")
	flag.StringVar(&(c.BaseAddress), "b", DefaultBaseAddress, "Base address before a shorted url")
	flag.StringVar(&(c.LogLevel), "l", DefaultLogLevel, "Log level")
	flag.StringVar(&(c.FileStoragePath), "f", DefaultFileStoragePath, "File storage path")
	flag.StringVar(&(c.DBConnString), "d", DefaultDBConnectionString, "DB connection string")
	flag.Parse()

	envServerAddress, wasFoundServerAddress := os.LookupEnv("SERVER_ADDRESS")
	envBaseAddress, wasFoundBaseAddress := os.LookupEnv("BASE_URL")
	envLogLevel, wasFoundLogLevel := os.LookupEnv("LOG_LEVEL")
	envFileStoragePath, wasFoundFileStoragePath := os.LookupEnv("FILE_STORAGE_PATH")
	envDBConnString, wasFoundDBConnString := os.LookupEnv("DATABASE_DSN")

	if c.ServerAddress == DefaultServerAddress && wasFoundServerAddress {
		c.ServerAddress = envServerAddress
	}
	if c.BaseAddress == DefaultBaseAddress && wasFoundBaseAddress {
		c.BaseAddress = envBaseAddress
	}
	if c.LogLevel == DefaultLogLevel && wasFoundLogLevel {
		c.LogLevel = envLogLevel
	}
	if wasFoundFileStoragePath {
		c.FileStoragePath = envFileStoragePath
	}
	if wasFoundDBConnString {
		c.DBConnString = envDBConnString
	}
	//`else` - flag value (it has been already set)
}
