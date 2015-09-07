// Copyright 2015 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package currency

import (
	"fmt"
	"testing"

	declib2 "github.com/shopspring/decimal"
	declib1 "github.com/wayn3h0/go-decimal"
)

func TestCurrencyConverter(t *testing.T) {
	//ast := assert.NewAssertWithName(t, "CurrencyConverter")

	cc := NewCurrencyConverter()
	cc.Update()
	fmt.Println(cc.ConvertString("1", USD, EUR))
	fmt.Println(cc.ConvertString("1", PLN, EUR))
	//fmt.Printf("%#v", Currency(EUR))
}

func TestDecimal1(t *testing.T) {
	dec1, err := declib1.Parse("0.1")
	if err != nil {
		t.Fatal(err)
	}

	dec2, err := declib1.Parse("0.2")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(dec1.Add(dec2))
	fmt.Println(dec1.Div(dec2))
}

func TestDecimal2(t *testing.T) {
	dec1, err := declib2.NewFromString("0.1")
	if err != nil {
		t.Fatal(err)
	}

	dec2, err := declib2.NewFromString("0.2")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(dec1.Add(dec2))
	fmt.Println(dec1.Div(dec2))
}

func TestDecimal3(t *testing.T) {
	one, _ := declib1.Parse("1.0")
	amount, _ := declib1.Parse("5.0")
	rate, _ := declib1.Parse("4.2")

	fmt.Println(amount.Mul(one.Div(rate)).FloatString(4))
}

func TestDecimal4(t *testing.T) {
	one, _ := declib2.NewFromString("1.0")
	amount, _ := declib2.NewFromString("5.0")
	rate, _ := declib2.NewFromString("4.2")

	fmt.Println(amount.Mul(one.Div(rate)).StringFixed(4))
}
