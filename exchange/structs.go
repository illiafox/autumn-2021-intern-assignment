package exchange

type Currency struct {
	Id  string  `json:"id"`
	Val float64 `json:"val"`
	To  string  `json:"to"`
	Fr  string  `json:"fr"`
}

type successJSON struct {
	Results map[string]Currency `json:"results"`
}
