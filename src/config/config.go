package config

const (
	DefaultConnectionString = "host=localhost port=5432 user=postgres password=postgres dbname=scurvy10k sslmode=disable"
)

type Config struct {
	ConnectionString string
}

type Option func(*Config)

func NewConfig(os ...Option) Config {
	c := &Config{
		ConnectionString: DefaultConnectionString,
	}
	for _, o := range os {
		o(c)
	}
	return *c
}

func WithConnectionString(s string) Option {
	return func(c *Config) {
		c.ConnectionString = s
	}
}
