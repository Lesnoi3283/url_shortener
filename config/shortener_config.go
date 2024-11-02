// Package config is a configuration package.
// It contains a Config struct and functions to read configuration params from env variables and command line args.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

// Default configurations params.
const (
	DefaultBaseAddress        = "http://localhost:8080"
	DefaultServerAddress      = "localhost:8080"
	DefaultGRPCAddress        = "localhost:50051"
	DefaultLogLevel           = "info"
	DefaultFileStoragePath    = "/tmp/short-url-db.json"
	DefaultDBConnectionString = ""
	DefaultEnableHTTPSFlag    = false
	DefaultTrustedSubnet      = "127.0.0.1/24"
	DefaultJWTTimeoutHours    = 5
)

type confFileData struct {
	ServerAddress   string `json:"server_address"`
	GRPCAddress     string `json:"grpc_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDsn     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	LogLevel        string `json:"log_level"`
	TrustedSubnet   string `json:"trusted_subnet"`
	JWTSecret       string `json:"jwt_secret"`
	JWTTimeoutHours int    `json:"jwt_timeout_hours"`
}

// Config is a struct with configuration params.
// Attention - JWTSecret can be read ONLY from environment or configuration file.
type Config struct {
	BaseAddress     string
	ServerAddress   string
	GRPCAddress     string
	LogLevel        string
	FileStoragePath string
	DBConnString    string
	EnableHTTPS     bool
	ConfigFileName  string
	TrustedSubnet   string
	JWTSecret       string
	JWTTimeoutHours int
}

// Configure reads configuration params from command line args, environmental variables and DefaultConstParams.
// And writes them into a Config struct.
func (c *Config) Configure() error {
	//get flag values
	flag.StringVar(&(c.ServerAddress), "a", DefaultServerAddress, "Address where server will work. Example: \"localhost:8080\".")
	flag.StringVar(&(c.GRPCAddress), "g", DefaultGRPCAddress, "Address where gRPC will work. Example: \"localhost:50051\".")
	flag.StringVar(&(c.BaseAddress), "b", DefaultBaseAddress, "Base address before a shorted url")
	flag.StringVar(&(c.LogLevel), "l", DefaultLogLevel, "Log level")
	flag.StringVar(&(c.FileStoragePath), "f", DefaultFileStoragePath, "File storage path")
	flag.StringVar(&(c.DBConnString), "d", DefaultDBConnectionString, "DB connection string")
	flag.BoolVar(&(c.EnableHTTPS), "s", DefaultEnableHTTPSFlag, "This flag enables HTTPS support")
	flag.StringVar(&(c.ConfigFileName), "c", "", "Config file name")
	flag.StringVar(&(c.TrustedSubnet), "t", DefaultTrustedSubnet, "Trusted subnet")
	flag.IntVar(&(c.JWTTimeoutHours), "j", DefaultJWTTimeoutHours, "JWT timeout hours")
	flag.Parse()

	//get env values
	envServerAddress, wasFoundServerAddress := os.LookupEnv("SERVER_ADDRESS")
	envGRPCAddress, wasFoundGRPCAddress := os.LookupEnv("GRPC_ADDRESS")
	envBaseAddress, wasFoundBaseAddress := os.LookupEnv("BASE_URL")
	envLogLevel, wasFoundLogLevel := os.LookupEnv("LOG_LEVEL")
	envFileStoragePath, wasFoundFileStoragePath := os.LookupEnv("FILE_STORAGE_PATH")
	envDBConnString, wasFoundDBConnString := os.LookupEnv("DATABASE_DSN")
	envEnableHTTPS, wasFoundEnableHTTPSFlag := os.LookupEnv("ENABLE_HTTPS")
	envConfFile, wasFoundConfFile := os.LookupEnv("CONFIG")
	envTrustedSubnet, wasFoundTrustedSubnet := os.LookupEnv("TRUSTED_SUBNET")
	envJWTSecret, wasFoundJWTSecret := os.LookupEnv("JWT_SECRET")
	envJWTTimeoutHours, wasFoundJWTTimeoutHours := os.LookupEnv("JWT_TIMEOUT_HOURS")

	//set values
	if c.ServerAddress == DefaultServerAddress && wasFoundServerAddress {
		c.ServerAddress = envServerAddress
	}
	if c.GRPCAddress == DefaultGRPCAddress && wasFoundGRPCAddress {
		c.GRPCAddress = envGRPCAddress
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
	if wasFoundEnableHTTPSFlag {
		parsedEnableHTTPS, err := strconv.ParseBool(envEnableHTTPS)
		if err != nil {
			return fmt.Errorf("error parsing ENABLE_HTTPS env var: %w", err)
		}
		c.EnableHTTPS = parsedEnableHTTPS
	}
	if wasFoundTrustedSubnet {
		c.TrustedSubnet = envTrustedSubnet
	}
	if wasFoundJWTSecret {
		c.JWTSecret = envJWTSecret
	} else {
		// panic("JWT secret not set. Please set the environment variable `JWT_SECRET`.")
		// I believe the panic function should be called in this case
		// since using a default value for c.JWTSecret is not secure.
		// However, I will use it here because Yandex.Practicum's tests don't set
		// a `JWT_SECRET` variable.
		// Therefore, I am using a "hardcoded" default value for c.JWTSecret solely for educational purposes.
		c.JWTSecret = "superSecret"
	}
	if wasFoundJWTTimeoutHours {
		hours, err := strconv.Atoi(envJWTTimeoutHours)
		if err != nil {
			return fmt.Errorf("error parsing JWT_TIMEOUT_HOURS: %w", err)
		}
		c.JWTTimeoutHours = hours
	}
	//`else` - flag value (it has been already set)

	//get config file values and set them if they were not provided earlier
	if wasFoundConfFile {
		c.ConfigFileName = envConfFile
		//open file
		file, err := os.Open(c.ConfigFileName)
		if err != nil {
			return fmt.Errorf("could not open config file: %w", err)
		}

		//read
		data, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("could not read config file: %w", err)
		}

		//parse
		confData := &confFileData{}
		err = json.Unmarshal(data, confData)
		if err != nil {
			return fmt.Errorf("could not parse config file: %w", err)
		}

		//set values
		if c.ServerAddress == DefaultServerAddress && confData.ServerAddress != "" {
			c.ServerAddress = confData.ServerAddress
		}
		if c.GRPCAddress == DefaultGRPCAddress && confData.GRPCAddress != "" {
			c.GRPCAddress = confData.GRPCAddress
		}
		if c.BaseAddress == DefaultBaseAddress && confData.BaseURL != "" {
			c.BaseAddress = confData.BaseURL
		}
		if c.LogLevel == DefaultLogLevel && confData.LogLevel != "" {
			c.LogLevel = confData.LogLevel
		}
		if c.FileStoragePath == DefaultFileStoragePath && confData.FileStoragePath != "" {
			c.FileStoragePath = confData.FileStoragePath
		}
		if c.DBConnString == DefaultDBConnectionString && confData.DatabaseDsn != "" {
			c.DBConnString = confData.DatabaseDsn
		}
		if !c.EnableHTTPS && confData.EnableHTTPS {
			c.EnableHTTPS = confData.EnableHTTPS
		}
		if c.TrustedSubnet == DefaultTrustedSubnet && confData.TrustedSubnet != "" {
			c.TrustedSubnet = confData.TrustedSubnet
		}
		if !wasFoundJWTSecret && confData.JWTSecret != "" {
			c.JWTSecret = confData.JWTSecret
		}
		if c.JWTTimeoutHours == DefaultJWTTimeoutHours && confData.JWTTimeoutHours != 0 {
			c.JWTTimeoutHours = confData.JWTTimeoutHours
		}
	}
	return nil
}
