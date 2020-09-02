package currency

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Build interface {
	BuildResponse(c *Currency)
}

type Response struct {
	Base  string
	Date  string
	Rates map[string]float64
}

type Config struct {
	Host string
}

type Currency struct {
	Config   *Config
	Response *Response
}

func (c *Currency) GetResponse(f interface{}) error {
	r, err := http.Get(c.Config.Host)
	if err != nil {
		return err
	}

	err = json.NewDecoder(r.Body).Decode(f)
	if err != nil {
		return err
	}

	return nil
}

func (c *Currency) getRate(symbol string) (float64, error) {
	r := c.Response.Rates[strings.ToUpper(symbol)]
	if r == 0 {
		return r, fmt.Errorf("%s Symbol not found!", symbol)
	}

	return r, nil
}

func (c *Currency) Exchange(p float64, symbol string) (float64, error) {
	r, err := c.getRate(symbol)
	if err != nil {
		return 0, err
	}

	return p * r, nil
}
