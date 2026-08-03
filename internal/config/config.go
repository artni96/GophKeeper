package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

type ConfigFile struct {
	ServerAddr string `json:"server_addr"`

	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBName     string `json:"db_name"`
	SSLMode    string `json:"ssl_mode"`

	SecretKey  string        `json:"secret_key"`
	TokenExp   time.Duration `json:"token_exp"`
	EnableTCP  bool          `json:"enable_tcp"`
	CertFile   string        `json:"cert_file"`
	KeyFile    string        `json:"key_file"`
	LoggingLvl string        `json:"logging_lvl"`

	GPeriod time.Duration `json:"grace_period"`
	FPeriod time.Duration `json:"force_period"`
}
type Config struct {
	ServerAddr string
	DBDsn      string
	DBName     string
	SecretKey  string
	TokenExp   time.Duration
	EnableTCP  bool
	CertFile   string
	KeyFile    string
	LoggingLvl string
	GPeriod    time.Duration
	FPeriod    time.Duration
}

func ParseConfig() (*Config, error) {
	conf := Config{}
	confFile := ConfigFile{}

	fs := flag.NewFlagSet("config", flag.ExitOnError)

	var configFilePath string
	declaredFlags := make(map[string]bool)

	fs.StringVar(&configFilePath, "c", "", "json config file")

	fs.StringVar(&conf.ServerAddr, "a", "", "server address")
	fs.BoolVar(&conf.EnableTCP, "t", false, "enable tps")
	fs.StringVar(&conf.CertFile, "cf", "", "cert file")
	fs.StringVar(&conf.KeyFile, "kf", "", "key file")
	fs.StringVar(&conf.LoggingLvl, "l", "", "logging level")

	var tokenExpMins int

	fs.IntVar(&tokenExpMins, "e", 0, "token expiration time")

	err := fs.Parse(os.Args[1:])

	conf.TokenExp = time.Duration(tokenExpMins) * time.Minute
	if err != nil {
		return &conf, fmt.Errorf("failed to parse flags for config: %w", err)
	}

	fs.Visit(func(f *flag.Flag) {
		declaredFlags[f.Name] = true
	})

	if configFilePath != "" {
		err = confFile.parseFile(configFilePath)
		if err != nil {
			return &conf, err
		}
	}

	var errs []string

	envServerAddr, ok := os.LookupEnv("SERVER_ADDR")
	if ok && !declaredFlags["a"] {
		conf.ServerAddr = envServerAddr
	} else if confFile.ServerAddr != "" {
		conf.ServerAddr = confFile.ServerAddr
	} else {
		errs = append(errs, "server address is not set")
	}

	envEnableTPS, ok := os.LookupEnv("ENABLE_TPS")
	if ok && !declaredFlags["t"] {
		conf.EnableTCP = envEnableTPS == "true"
	} else if confFile.EnableTCP {
		conf.EnableTCP = confFile.EnableTCP
	}

	dbDSN, err := confFile.AssembleDBDsn()
	if err != nil {
		errs = append(errs, err.Error())
	} else {
		conf.DBDsn = dbDSN
		conf.DBName = confFile.DBName
	}

	envTokenExp, ok := os.LookupEnv("TOKEN_EXP")
	if ok && !declaredFlags["e"] {
		t, err := strconv.Atoi(envTokenExp)
		if err != nil {
			errs = append(errs, err.Error())
		}
		conf.TokenExp = time.Duration(t) * time.Minute
	} else if confFile.TokenExp != 0 {
		conf.TokenExp = confFile.TokenExp * time.Minute
	} else if !declaredFlags["e"] {
		errs = append(errs, "token expiration is not set")
	}

	envCertFile, ok := os.LookupEnv("CERT_FILE")
	if ok && !declaredFlags["cf"] {
		conf.CertFile = envCertFile
	} else if confFile.CertFile != "" {
		conf.CertFile = confFile.CertFile
	} else if conf.EnableTCP {
		errs = append(errs, "certificate file is not set")
	}

	envKeyFile, ok := os.LookupEnv("KEY_FILE")
	if ok && !declaredFlags["kf"] {
		conf.KeyFile = envKeyFile
	} else if confFile.KeyFile != "" {
		conf.KeyFile = confFile.KeyFile
	} else if conf.EnableTCP {
		errs = append(errs, "key file is not set")
	}

	envLoggingLvl, ok := os.LookupEnv("LOGGING_LVL")
	if ok && !declaredFlags["l"] {
		conf.LoggingLvl = envLoggingLvl
	} else if confFile.LoggingLvl != "" {
		conf.LoggingLvl = confFile.LoggingLvl
	} else {
		conf.LoggingLvl = zapcore.DebugLevel.String()
	}

	if confFile.SecretKey != "" {
		conf.SecretKey = confFile.SecretKey
	} else {
		errs = append(errs, "secret key is not set")
	}

	if confFile.GPeriod != 0 {
		conf.GPeriod = confFile.GPeriod
	} else {
		errs = append(errs, "grace period is not set")
	}

	if confFile.FPeriod != 0 {
		conf.FPeriod = confFile.FPeriod
	} else {
		errs = append(errs, "force period is not set")
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}
	return &conf, nil
}

// parseFile collects data from the file filename into ConfigFile.
func (f *ConfigFile) parseFile(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	err = json.Unmarshal(file, &f)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	return nil
}

// AssembleDBDsn parses app config and returns a database destination string.
func (f *ConfigFile) AssembleDBDsn() (string, error) {
	var errs []string
	if f.DBHost == "" {
		errs = append(errs, "db host is not set")
	}
	if f.DBPort == 0 {
		errs = append(errs, "db port is not set")
	}
	if f.DBUser == "" {
		errs = append(errs, "db user is not set")
	}
	if f.DBPassword == "" {
		errs = append(errs, "db password is not set")
	}
	if f.DBName == "" {
		errs = append(errs, "db name is not set")
	}
	if f.SSLMode == "" {
		f.SSLMode = "disable"
	}
	if len(errs) > 0 {
		return "", errors.New("failed to parse database destination: " + strings.Join(errs, ", "))
	}

	dbDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		f.DBHost, f.DBPort, f.DBUser, f.DBPassword, f.DBName, f.SSLMode)
	return dbDSN, nil
}
