package exchange

import (
	"autumn-2021-intern-assignment/utils/config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Update(conf config.Exchanger) error {
	u, err := url.Parse(conf.Endpoint)
	if err != nil {
		return fmt.Errorf("parsing endpoint: %w", err)
	}

	query := u.Query()
	query.Set("apiKey", conf.Key)

	if len(conf.Bases) != 0 {
		for i := range conf.Bases {
			conf.Bases[i] += "_RUB"
		}
		query.Set("q", strings.Join(conf.Bases, ","))
	}

	resp, err := http.Get(u.String() + "?" + query.Encode())
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}

	if resp.StatusCode != 200 {
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(data))
	}

	var result successJSON
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("decoding json: %w", err)
	}

	for _, c := range result.Results {
		exchangeMap.Store(c.Fr, c.Val)
	}

	return nil
}
