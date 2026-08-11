package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap/zapcore"
)

// Config represents main setting for the server.
type Config struct {
	ServerAddr string
	DBDsn      string
	DBName     string
	SecretKey  string
	TokenExp   time.Duration
	EnableTCP  bool
	CertFile   string
	KeyFile    string
	LogLvl     string
	GPeriod    time.Duration
	FPeriod    time.Duration
}

// ParseConfig parse flags and environment variables and applies them to the server config.
func ParseConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}
	conf := Config{}

	fs := flag.NewFlagSet("config", flag.ExitOnError)

	declaredFlags := make(map[string]bool)

	var servHost string
	var servPort int

	fs.StringVar(&servHost, "h", "", "server host")
	fs.IntVar(&servPort, "p", 0, "server port")

	fs.BoolVar(&conf.EnableTCP, "t", false, "enable tps")
	fs.StringVar(&conf.CertFile, "cf", "", "cert file")
	fs.StringVar(&conf.KeyFile, "kf", "", "key file")
	fs.StringVar(&conf.LogLvl, "l", "", "logging level")
	fs.StringVar(&conf.SecretKey, "s", "", "server key")

	var dbName, dbHost, dbUser, dbPass, sslMode string
	var dbPort int

	fs.StringVar(&dbName, "db", "", "database name")
	fs.StringVar(&dbHost, "dbh", "", "database host")
	fs.IntVar(&dbPort, "dbp", 0, "database port")
	fs.StringVar(&dbUser, "dbu", "", "database user")
	fs.StringVar(&dbPass, "dbpw", "", "database password")
	fs.StringVar(&sslMode, "dbs", "disable", "database ssl mode")

	var tokenExpMins int

	fs.IntVar(&tokenExpMins, "e", 0, "token expiration time")

	var gracePeriodMins int
	var forcePeriodMins int

	fs.IntVar(&gracePeriodMins, "gp", 0, "grace period mins")
	fs.IntVar(&forcePeriodMins, "fp", 0, "force period mins")

	err = fs.Parse(os.Args[1:])
	if err != nil {
		return &conf, fmt.Errorf("failed to parse flags for config: %w", err)
	}

	conf.TokenExp = time.Duration(tokenExpMins) * time.Minute
	conf.GPeriod = time.Duration(gracePeriodMins) * time.Minute
	conf.FPeriod = time.Duration(forcePeriodMins) * time.Minute

	fs.Visit(func(f *flag.Flag) {
		declaredFlags[f.Name] = true
	})

	var errs []string

	if !declaredFlags["h"] {
		envServHost, ok := os.LookupEnv("SERVER_HOST")
		if ok {
			servHost = envServHost
		} else {
			errs = append(errs, "server host is not set")
		}
	}

	if !declaredFlags["p"] {
		envServPort, ok := os.LookupEnv("SERVER_PORT")
		if ok {
			servPort, err = strconv.Atoi(envServPort)
			if err != nil {
				errs = append(errs, "invalid server port value")
			}
		} else {
			errs = append(errs, "server port is not set")
		}
	}

	if servHost != "" && servPort != 0 {
		conf.ServerAddr = fmt.Sprintf("%s:%d", servHost, servPort)
	}

	if !declaredFlags["db"] {
		envDBName, ok := os.LookupEnv("DB_NAME")
		if ok {
			dbName = envDBName
			conf.DBName = dbName
		} else {
			errs = append(errs, "db name not set")
		}
	}

	if !declaredFlags["dbh"] {
		envDBHost, ok := os.LookupEnv("DB_HOST")
		if ok {
			dbHost = envDBHost
		} else {
			errs = append(errs, "database host is not set")
		}
	}

	if !declaredFlags["dbp"] {
		envDBPort, ok := os.LookupEnv("DB_PORT")
		if ok {
			dbPort, err = strconv.Atoi(envDBPort)
			if err != nil {
				errs = append(errs, "invalid database port")
			}
		} else {
			errs = append(errs, "database port is not set")
		}
	}

	if !declaredFlags["dbu"] {
		envDBUser, ok := os.LookupEnv("DB_USER")
		if ok {
			dbUser = envDBUser
		} else {
			errs = append(errs, "database user is not set")
		}
	}

	if !declaredFlags["dbpw"] {
		envDBPassword, ok := os.LookupEnv("DB_PASSWORD")
		if ok {
			dbPass = envDBPassword
		} else {
			errs = append(errs, "database password is not set")
		}
	}

	if !declaredFlags["dbs"] {
		envDBSSLMode, ok := os.LookupEnv("SSL_MODE")
		if ok {
			if envDBSSLMode != "enable" && envDBSSLMode != "disable" {
				errs = append(errs, "invalid ssl mode value")
			} else {
				sslMode = envDBSSLMode
			}
		} else {
			errs = append(errs, "database ssl mode is not set")
		}
	}

	if dbName != "" && dbHost != "" && dbUser != "" && dbPass != "" && dbPort != 0 && sslMode != "" {
		conf.DBDsn = assembleDBDsn(dbName, dbHost, dbUser, dbPass, sslMode, dbPort)
	}

	if !declaredFlags["t"] {
		envEnableTCP, ok := os.LookupEnv("ENABLE_TCP")
		if ok {
			if envEnableTCP != "true" && envEnableTCP != "false" {
				errs = append(errs, "invalid enable tcp value")
			} else {
				conf.EnableTCP = envEnableTCP == "true"
			}
		} else {
			errs = append(errs, "enable tcp port is not set")
		}
	}

	if !declaredFlags["e"] {
		envTokenExp, ok := os.LookupEnv("TOKEN_EXP")
		if ok {
			t, err := strconv.Atoi(envTokenExp)
			if err != nil {
				errs = append(errs, "invalid token expiration time")
			} else {
				conf.TokenExp = time.Duration(t) * time.Minute
			}
		} else {
			errs = append(errs, "token expiration time is not set")
		}
	}

	if !declaredFlags["cf"] {
		envCertFile, ok := os.LookupEnv("CERT_FILE")
		if ok {
			conf.CertFile = envCertFile
		} else if conf.EnableTCP {
			errs = append(errs, "cert file is not set")
		}
	}

	if !declaredFlags["kf"] {
		envKeyFile, ok := os.LookupEnv("KEY_FILE")
		if ok {
			conf.KeyFile = envKeyFile
		} else if conf.EnableTCP {
			errs = append(errs, "key file is not set")
		}
	}

	if !declaredFlags["l"] {
		envLogLevel, ok := os.LookupEnv("LOG_LEVEL")
		if ok {
			conf.LogLvl = envLogLevel
		} else {
			conf.LogLvl = zapcore.DebugLevel.String()
		}
	}

	if !declaredFlags["s"] {
		envSecretKey, ok := os.LookupEnv("SECRET_KEY")
		if ok {
			conf.SecretKey = envSecretKey
		} else {
			errs = append(errs, "secret key is not set")
		}
	}

	if !declaredFlags["gp"] {
		envGracePeriod, ok := os.LookupEnv("GRACE_PERIOD")
		if ok {
			t, err := strconv.Atoi(envGracePeriod)
			if err != nil {
				errs = append(errs, "invalid grace period value")
			} else {
				conf.GPeriod = time.Duration(t) * time.Second
			}
		} else {
			errs = append(errs, "grace period time is not set")
		}
	}

	if !declaredFlags["fp"] {
		envFPeriod, ok := os.LookupEnv("FORCE_PERIOD")
		if ok {
			t, err := strconv.Atoi(envFPeriod)
			if err != nil {
				errs = append(errs, "invalid force period value")
			} else {
				conf.FPeriod = time.Duration(t) * time.Second
			}
		} else {
			errs = append(errs, "force period time is not set")
		}
	}

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}
	return &conf, nil
}

func assembleDBDsn(dbName, dbHost, dbUser, dbPassword, sslMode string, dbPort int) string {
	dbDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
	return dbDSN
}
