package exchange

import (
	"autumn-2021-intern-assignment/utils/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func Update(conf config.Exchanger) error {
	u, err := url.Parse(conf.Endpoint)
	if err != nil {
		return fmt.Errorf("parsing endpoint: %w", err)
	}

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	if len(conf.Base) == 0 {
		resp, err := http.Get(u.String() + "currencies?apiKey=" + conf.Key)
		if err != nil {
			return fmt.Errorf("send currencies request: %w", err)
		}

		var curr = struct {
			Error   string `json:"error"`
			Results map[string]struct {
				ID string `json:"id"`
			} `json:"results"`
		}{}

		err = json.NewDecoder(resp.Body).Decode(&curr)

		if err != nil {
			return fmt.Errorf("decoding currencies json: %w", err)
		}

		if len(curr.Results) == 0 {
			return fmt.Errorf("no currencies results: %s", curr.Error)
		}

		conf.Base = make([]string, 0, len(curr.Results))

		for k := range curr.Results {
			conf.Base = append(conf.Base, k+"_RUB")
		}
	}

	for i := range conf.Base {
		conf.Base[i] += "_RUB"
	}

	query := u.Query()
	query.Set("apiKey", conf.Key)
	query.Set("compact", "ultra")

	for i := 2; i <= len(conf.Base); i += 2 {
		query.Set("q", strings.Join(conf.Base[i-2:i], ","))
		resp, err := client.Get(u.String() + "convert?" + query.Encode())
		if err != nil {
			return fmt.Errorf("send exchange request: %w", err)
		}

		if resp.StatusCode != 200 {
			data, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(data))
		}

		var result = make(map[string]float64)
		err = json.NewDecoder(resp.Body).Decode(&result)

		if err != nil {
			return fmt.Errorf("decoding json: %w", err)
		}

		for k, v := range result {
			exchangeMap.Store(strings.TrimSuffix(k, "_RUB"), v)
		}
	}

	return nil
}

func Store(filepath string) error {
	exchanges := make(map[string]float64)

	exchangeMap.Range(func(key, value interface{}) bool {
		attr, ok := key.(string)
		if !ok {
			return false
		}
		curr, ok := value.(float64)
		if !ok {
			return false
		}

		exchanges[attr] = curr

		return true
	})

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	err = json.NewEncoder(file).Encode(exchanges)
	if err != nil {
		return fmt.Errorf("encoding map: %w", err)
	}

	return nil
}

func Load(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	exchanges := make(map[string]float64)

	err = json.Unmarshal(data, &exchanges)
	if err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}

	for k, v := range exchanges {
		exchangeMap.Store(k, v)
	}

	return nil
}
