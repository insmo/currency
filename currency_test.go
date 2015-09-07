// Copyright 2015 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package currency

import "testing"

func TestCurrencyConverter(t *testing.T) {
	cc := New()
	usdtoeur, err := cc.ConvertString("1", USD, EUR)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("1 USD to EUR = %s", usdtoeur.StringFixed(4))

	plntoeur, err := cc.ConvertString("1", PLN, EUR)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("1 PLN to EUR = %s", plntoeur.StringFixed(4))
}
