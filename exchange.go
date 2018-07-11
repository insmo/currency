package currency

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

var oneD = decimal.NewFromFloat(1.0)

type date string

func toDate(t time.Time) date {
	return date(t.Format("20060102"))
}

func toFixerDate(t time.Time) date {
	return date(t.Format("2006-01-02"))
}

type ExchangeRate struct {
	FromEUR decimal.Decimal
	ToEUR   decimal.Decimal
}

// Exchange holds a cache of currency exchange rates.
type Exchange struct {
	cache    map[date]map[Currency]ExchangeRate
	mux      sync.Mutex
	apiToken string
}

// NewExchange initializes a new Exchange.
func NewExchange(apiToken string) *Exchange {
	return &Exchange{
		cache:    make(map[date]map[Currency]ExchangeRate),
		apiToken: apiToken,
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
	fixerData, err := fetchFixerData(t, ex.apiToken)

	if err != nil {
		return err
	}

	data, err := normalizeFixerData(fixerData)

	if err != nil {
		return err
	}

	ex.cache[toDate(t)] = data
	return nil
}

type fixerCurrencyResponse struct {
	Base    string `json:"base"`
	Date    string `json:"date"`
	Success bool   `json:"success"`
	Error   struct {
		Info string `json:"info"`
	} `json:"error,omitempty"`
	//Rates interface{} `json:"rate"`
	Rates struct {
		AUD float64 `json:"AUD"`
		BGN float64 `json:"BGN"`
		BRL float64 `json:"BRL"`
		CAD float64 `json:"CAD"`
		CHF float64 `json:"CHF"`
		CNY float64 `json:"CNY"`
		CYP float64 `json:"CYP"`
		CZK float64 `json:"CZK"`
		DKK float64 `json:"DKK"`
		EEK float64 `json:"EEK"`
		EUR float64 `json:"EUR"`
		GBP float64 `json:"GBP"`
		HKD float64 `json:"HKD"`
		HRK float64 `json:"HRK"`
		HUF float64 `json:"HUF"`
		IDR float64 `json:"IDR"`
		ILS float64 `json:"ILS"`
		INR float64 `json:"INR"`
		ISK float64 `json:"ISK"`
		JPY float64 `json:"JPY"`
		KRW float64 `json:"KRW"`
		LTL float64 `json:"LTL"`
		LVL float64 `json:"LVL"`
		MTL float64 `json:"MTL"`
		MXN float64 `json:"MXN"`
		MYR float64 `json:"MYR"`
		NOK float64 `json:"NOK"`
		NZD float64 `json:"NZD"`
		PHP float64 `json:"PHP"`
		PLN float64 `json:"PLN"`
		ROL float64 `json:"ROL"`
		RON float64 `json:"RON"`
		RUB float64 `json:"RUB"`
		SEK float64 `json:"SEK"`
		SGD float64 `json:"SGD"`
		SIT float64 `json:"SIT"`
		SKK float64 `json:"SKK"`
		THB float64 `json:"THB"`
		TRL float64 `json:"TRL"`
		TRY float64 `json:"TRY"`
		USD float64 `json:"USD"`
		ZAR float64 `json:"ZAR"`
	} `json:"rates"`
}

func normalizeFixerData(fixerData *fixerCurrencyResponse) (map[Currency]ExchangeRate, error) {
	data := make(map[Currency]ExchangeRate)

	add := func(cur Currency, price float64) {
		fromEUR := decimal.NewFromFloat(price)
		toEUR := fromEUR

		if price != 0 {
			toEUR = oneD.Div(fromEUR)
		}

		data[cur] = ExchangeRate{
			FromEUR: fromEUR,
			ToEUR:   toEUR,
		}
	}

	add(AUD, fixerData.Rates.AUD)
	add(BGN, fixerData.Rates.BGN)
	add(BRL, fixerData.Rates.BRL)
	add(CAD, fixerData.Rates.CAD)
	add(CHF, fixerData.Rates.CHF)
	add(CNY, fixerData.Rates.CNY)
	add(CYP, fixerData.Rates.CYP)
	add(CZK, fixerData.Rates.CZK)
	add(DKK, fixerData.Rates.DKK)
	//add(EUR, fixerData.Rates.EUR)
	add(GBP, fixerData.Rates.GBP)
	add(HKD, fixerData.Rates.HKD)
	add(HRK, fixerData.Rates.HRK)
	add(HUF, fixerData.Rates.HUF)
	add(IDR, fixerData.Rates.IDR)
	add(ILS, fixerData.Rates.ILS)
	add(INR, fixerData.Rates.INR)
	add(ISK, fixerData.Rates.ISK)
	add(JPY, fixerData.Rates.JPY)
	add(KRW, fixerData.Rates.KRW)
	add(LTL, fixerData.Rates.LTL)
	add(LVL, fixerData.Rates.LVL)
	add(MXN, fixerData.Rates.MXN)
	add(MYR, fixerData.Rates.MYR)
	add(NOK, fixerData.Rates.NOK)
	add(NZD, fixerData.Rates.NZD)
	add(PHP, fixerData.Rates.PHP)
	add(PLN, fixerData.Rates.PLN)
	add(RON, fixerData.Rates.RON)
	add(RUB, fixerData.Rates.RUB)
	add(SEK, fixerData.Rates.SEK)
	add(SGD, fixerData.Rates.SGD)
	add(SIT, fixerData.Rates.SIT)
	add(THB, fixerData.Rates.THB)
	add(TRY, fixerData.Rates.TRY)
	add(USD, fixerData.Rates.USD)
	add(ZAR, fixerData.Rates.ZAR)

	data[EUR] = ExchangeRate{
		FromEUR: decimal.NewFromFloat(1.0),
		ToEUR:   decimal.NewFromFloat(1.0),
	}

	return data, nil
}

func fetchFixerData(t time.Time, apiToken string) (*fixerCurrencyResponse, error) {
	maxTries := 1

	for i := 0; i < maxTries; i++ {
		resp, err := fixerDataRequest(t, apiToken)

		if err != nil {
			if i+1 == maxTries {
				return nil, err
			}
		} else {
			return resp, nil
		}

		if i == maxTries {
			time.Sleep(time.Millisecond * 100 * time.Duration(i))
		}
	}

	return nil, ErrFetchingData
}

func fixerDataRequest(t time.Time, apiToken string) (*fixerCurrencyResponse, error) {
	url := "http://data.fixer.io/api/" + string(toFixerDate(t)) + "?base=EUR&access_key=" + apiToken
	r, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	target := new(fixerCurrencyResponse)
	err = dec.Decode(target)

	if err != nil {
		return nil, err
	}

	if !target.Success {
		return nil, fmt.Errorf("fixer API err: %s", target.Error.Info)
	}

	return target, nil
}
