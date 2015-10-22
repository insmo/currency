package currency

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

var oneD = decimal.NewFromFloat(1.0)

type date string

func toDate(t time.Time) date {
	return date(t.Format("20060102"))
}

type ExchangeRate struct {
	FromUSD decimal.Decimal
	ToUSD   decimal.Decimal
}

// Exchange holds a cache of currency exchange rates.
type Exchange struct {
	cache map[date]map[Currency]ExchangeRate
	mux   sync.Mutex
}

// NewExchange initializes a new Exchange.
func NewExchange() *Exchange {
	return &Exchange{
		cache: make(map[date]map[Currency]ExchangeRate),
	}
}

// Get looks for an exchange rate for a given currency and date. It will update
// the cache it does not contain the exchange rates for the given date. It
// returns ErrNotExist if the exchange rate for the currency does not exist.
// It's safe to call Get concurrently from multiple go routines.
func (ex *Exchange) Get(t time.Time, c Currency) (ExchangeRate, error) {
	ex.mux.Lock()
	defer ex.mux.Unlock()
	key := toDate(t)
	day, ok := ex.cache[key]

	if !ok {
		err := ex.update(t)

		if err != nil {
			return ExchangeRate{}, err
		}

		day = ex.cache[key]
	}

	rate, ok := day[c]

	if !ok {
		return ExchangeRate{}, ErrNotExist{Currency: c, Time: t}
	}

	return rate, nil
}

func (ex *Exchange) update(t time.Time) error {
	yahooData, err := ex.fetchCurrencyData(t)

	if err != nil {
		return err
	}

	data, err := ex.normalizeCurrencyData(yahooData)

	if err != nil {
		return err
	}

	ex.cache[toDate(t)] = data
	return nil
}

func (ex *Exchange) normalizeCurrencyData(yahooData *yahooCurrencyResponse) (map[Currency]ExchangeRate, error) {
	data := make(map[Currency]ExchangeRate)

	for _, res := range yahooData.List.Resources {
		sym := res.Resource.Fields.Symbol

		// exp EUR=X
		if len(sym) != 5 {
			return nil, ErrCurrencyLength
		}

		cur, err := ParseCurrency(sym[:3])

		if err != nil {
			if err == ErrCurrencyUnknown {
				continue
			}
			return nil, err
		}

		price := res.Resource.Fields.Price

		// extra check
		if _, err := strconv.ParseFloat(price, 64); err != nil {
			return nil, err
		}

		fromUSD, err := decimal.NewFromString(price)

		if err != nil {
			return nil, err
		}

		data[cur] = ExchangeRate{
			FromUSD: fromUSD,
			ToUSD:   oneD.Div(fromUSD),
		}
	}

	return data, nil
}

type yahooCurrencyResponse struct {
	List struct {
		Meta struct {
			Count int    `json:"count"`
			Start int    `json:"start"`
			Type  string `json:"type"`
		} `json:"meta"`
		Resources []struct {
			Resource struct {
				Classname string `json:"classname"`
				Fields    struct {
					Date   string `json:"date"`
					Price  string `json:"price"`
					Symbol string `json:"symbol"`
					Type   string `json:"type"`
				} `json:"fields"`
			} `json:"resource"`
		} `json:"resources"`
	} `json:"list"`
}

func (ex *Exchange) fetchCurrencyData(t time.Time) (*yahooCurrencyResponse, error) {
	maxTries := 7

	for i := 1; i < maxTries; i++ {
		resp, err := fetchYahooData(t)

		if err != nil {
			if i == maxTries {
				return nil, err
			}
		} else if len(resp.List.Resources) > 0 {
			return resp, nil
		}

		if i == maxTries {
			time.Sleep(time.Millisecond * 100 * time.Duration(i))
		}
	}

	return nil, ErrFetchingData
}

func fetchYahooData(t time.Time) (*yahooCurrencyResponse, error) {
	r, err := http.Get("http://finance.yahoo.com/connection/currency-converter-cache?date=" + string(toDate(t)))

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, err
	}

	// Yahoo finance cache returns JavaScript. Remove that.
	body = bytes.TrimSpace(body)
	body = bytes.TrimLeft(body, "/**/YAHOO.Finance.CurrencyConverter.addConversionRates(")
	body = bytes.TrimRight(body, ");")

	target := new(yahooCurrencyResponse)
	err = json.Unmarshal(body, target)

	if err != nil {
		return nil, err
	}

	return target, nil
}
