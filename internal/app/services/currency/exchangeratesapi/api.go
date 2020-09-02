package exchangeratesapi

import (
	"github.com/creatortsv/finance-telegram-bot/internal/app/services/currency"
)

type APIResponse struct {
	currency.Build
	Base  string             `json:"base"`
	Date  string             `json:"date"`
	Rates map[string]float64 `json:"rates"`
}

func (r *APIResponse) BuildResponse(c *currency.Currency) {
	c.Response = &currency.Response{
		Base:  r.Base,
		Date:  r.Date,
		Rates: r.Rates,
	}
}

func New(base string) (*currency.Currency, error) {
	response := &APIResponse{}
	currency := &currency.Currency{
		Config: &currency.Config{
			Host: "https://api.exchangeratesapi.io/latest?base" + base,
		},
	}

	if err := currency.GetResponse(response); err != nil {
		return currency, err
	}

	response.BuildResponse(currency)

	return currency, nil
}
