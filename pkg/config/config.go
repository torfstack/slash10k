package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	DefaultHost     = "localhost"
	DefaultPort     = 5432
	DefaultUser     = "postgres"
	DefaultPassword = "postgres"
	DefaultDatabase = "slash10kdev"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Option func(*Config)

func NewConfig(os ...Option) Config {
	c := &Config{
		Host:     DefaultHost,
		Port:     DefaultPort,
		User:     DefaultUser,
		Password: DefaultPassword,
		Database: DefaultDatabase,
	}
	for _, o := range os {
		o(c)
	}
	return *c
}

func NewConfigFromEnv() (Config, error) {
	host := os.Getenv("DATABASE_CONNECTION_HOST")
	portS := os.Getenv("DATABASE_CONNECTION_PORT")
	port, err := strconv.Atoi(portS)
	user := os.Getenv("DATABASE_CONNECTION_USER")
	password := os.Getenv("DATABASE_CONNECTION_PASSWORD")
	database := os.Getenv("DATABASE_CONNECTION_DBNAME")

	if host == "" || portS == "" || err != nil || user == "" || password == "" || database == "" {
		return Config{}, errors.New("missing or malformed environment variables for database connection")
	}

	return NewConfig(
		WithHostName(host),
		WithPort(port),
		WithUser(user),
		WithPassword(password),
		WithDatabase(database),
	), nil
}

func (c Config) ConnectionString() string {
	return "host=" + c.Host + " port=" + strconv.Itoa(c.Port) + " user=" + c.User + " password=" + c.Password + " dbname=" + c.Database + " sslmode=disable"
}

func WithHostName(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}

func WithPort(port int) Option {
	return func(c *Config) {
		c.Port = port
	}
}

func WithUser(user string) Option {
	return func(c *Config) {
		c.User = user
	}
}

func WithPassword(password string) Option {
	return func(c *Config) {
		c.Password = password
	}
}

func WithDatabase(database string) Option {
	return func(c *Config) {
		c.Database = database
	}
}
