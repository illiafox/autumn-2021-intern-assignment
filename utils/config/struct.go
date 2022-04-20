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
}

type Exchanger struct {
	Skip     bool
	Every    int
	Endpoint string
	Key      string
	Base     []string
}

type Config struct {
	Postgres  Postgres
	Host      Host
	Exchanger Exchanger
}
