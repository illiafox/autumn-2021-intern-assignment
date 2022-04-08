package config

type Postgres struct {
	User     string
	Pass     string
	DbName   string
	IP       string
	Port     string
	Protocol string
}

type Host struct {
	Port string
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
