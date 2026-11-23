package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	AppEnvTest     = "test"
	appEnvKey      = "APP_ENV"
	projectName    = "buckwheat-queue"
	projectDirName = "buckwheat-queue"
)

var appEnv string

type Config struct {
	App ConfigApp
	DB  ConfigDB
}

type ConfigApp struct {
	Name     string
	Env      string
	LogLevel string
}

type ConfigDB struct {
	DBURL           string
	DBMigrationURL  string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func Load() (Config, error) {
	err := setEnv()
	if err != nil {
		return Config{}, err
	}

	logLevel, err := getEnv("LOG_LEVEL")
	dbHost, err := getEnv("DB_HOST")
	dbPort, err := getEnv("DB_PORT")
	dbName, err := getEnv("DB_NAME")
	dbUser, err := getEnv("DB_USER")
	dbPassword, err := getEnv("DB_PASSWORD")
	dbUserMigration, err := getEnv("DB_USER_MIGRATION")
	dbPasswordMigration, err := getEnv("DB_PASSWORD_MIGRATION")
	maxConns, err := getEnvInt("DB_MAX_CONNS")
	minConns, err := getEnvInt("DB_MIN_CONNS")
	maxConnLifetime, err := getEnvInt("DB_MAX_CONN_LIFETIME")
	maxConnIdleTime, err := getEnvInt("DB_MAX_CONN_IDLE_TIME")
	if err != nil {
		return Config{}, err
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)
	dbURLMigration := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUserMigration, dbPasswordMigration, dbHost, dbPort, dbName,
	)

	cfg := Config{
		App: ConfigApp{
			Name:     projectName,
			Env:      appEnv,
			LogLevel: logLevel,
		},
		DB: ConfigDB{
			DBURL:           dbURL,
			DBMigrationURL:  dbURLMigration,
			MaxConns:        int32(maxConns),
			MinConns:        int32(minConns),
			MaxConnLifetime: time.Duration(maxConnLifetime) * time.Second,
			MaxConnIdleTime: time.Duration(maxConnIdleTime) * time.Second,
		},
	}

	return cfg, nil
}

func setEnv() error {
	appEnv = os.Getenv(appEnvKey)
	if appEnv == "" || appEnv == AppEnvTest {
		re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
		currentPath, _ := os.Getwd()
		rootPath := string(re.Find([]byte(currentPath)))

		filenames := []string{rootPath + "/.env.dev.local", rootPath + "/.env.local", rootPath + "/.env"}
		if appEnv == AppEnvTest {
			filenames[0] = rootPath + "/.env.test.local"
		}

		err := godotenv.Load(filenames...)
		if err != nil {
			return fmt.Errorf("Loading .env files. Error: %w", err)
		}
		appEnv = os.Getenv(appEnvKey)
	}

	return nil
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", fmt.Errorf("Env value empty. Key: %s", key)
	}

	return value, nil
}

func getEnvInt(key string) (int, error) {
	value, err := getEnv(key)
	if err != nil {
		return 0, err
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("Error parse int from .env value. Key: %s. Error: %w", key, err)
	}

	return valueInt, nil
}
