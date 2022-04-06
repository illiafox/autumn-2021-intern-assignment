package exchange

type Currency struct {
	ID  string  `json:"id"`
	Val float64 `json:"val"`
	To  string  `json:"to"`
	Fr  string  `json:"fr"`
}

type successJSON struct {
	Results map[string]Currency `json:"results"`
}
