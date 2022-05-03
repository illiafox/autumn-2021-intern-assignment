package config

type Postgres struct {
	User     string `env:"POSTGRES_USER"`
	Pass     string `env:"POSTGRES_PASSWORD"`
	DbName   string `env:"POSTGRES_DATABASE"`
	IP       string `env:"POSTGRES_IP"`
	Port     string `env:"POSTGRES_PORT"`
	Protocol string `env:"POSTGRES_PROTOCOL"`
}

type Host struct {
	Port string `env:"HOST_PORT"`
	Key  string `env:"HOST_KEY_PATH"`  // Path to TLS key
	Cert string `env:"HOST_CERT_PATH"` // Path to TLS certificate
}

type Exchanger struct {
	Skip     bool   `env:"EXCHANGER_SKIP"`
	Load     bool   `env:"EXCHANGER_LOAD"`
	Every    int    `env:"EXCHANGER_EVERY"`
	Endpoint string `env:"EXCHANGER_ENDPOINT"`
	Key      string `env:"EXCHANGER_KEY"`
	Path     string `env:"EXCHANGER_PATH"`
	Base     []string
}

type Config struct {
	Postgres  Postgres
	Host      Host
	Exchanger Exchanger
}

func (c *Config) LoadEnv() error {
	return ReadEnv(c)
}
