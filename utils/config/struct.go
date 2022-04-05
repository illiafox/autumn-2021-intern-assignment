package config

type MySQL struct {
	Login    string
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
	Bases    []string
}

type Config struct {
	MySQL     MySQL
	Host      Host
	Exchanger Exchanger
}
