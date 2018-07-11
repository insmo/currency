// Copyright 2018 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package currency

import (
	"os"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestCurrencyConverter(t *testing.T) {
	token := os.Getenv("FIXER_API_TOKEN")

	if token == "" {
		t.Skip("FIXER_API_TOKEN required to run test")
	}

	cc := New(token)
	at := time.Date(2016, 9, 6, 0, 1, 0, 0, time.UTC)

	tests := []struct {
		value string
		from  Currency
		to    Currency
		exp   string
	}{
		{"1.0000", USD, USD, "1.0000"},
		{"1.0000", USD, EUR, "0.8884"},
		{"0.8884", EUR, USD, "1.0000"},
		{"1.0000", EUR, USD, "1.1256"},
		{"1.0000", PLN, USD, "0.2598"},
		{"1.0000", PLN, EUR, "0.2308"},
		{"0", PLN, EUR, "0.0000"},
	}

	for i, test := range tests {
		res, err := cc.ConvertStringAt(test.value, test.from, test.to, at)

		if err != nil {
			t.Fatal(err)
		}

		if res.StringFixed(4) != test.exp {
			t.Fatalf("test %d: expect %s, got %s", i, test.exp, res.StringFixed(4))
		}
	}
}

func TestCurrencyDivisionByZero(t *testing.T) {
	token := os.Getenv("FIXER_API_TOKEN")

	if token == "" {
		t.Skip("FIXER_API_TOKEN required to run test")
	}

	fromUSD, err := decimal.NewFromString("0")

	if err != nil {
		t.Fatal(err)
	}

	f, _ := fromUSD.Float64()

	if f != 0 {
		oneD.Div(fromUSD)
	}
}
