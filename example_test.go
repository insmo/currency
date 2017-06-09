package currency_test

import (
	"fmt"
	"log"
	"time"

	"github.com/simonz05/currency"
)

func ExampleConvertStringAt() {
	at := time.Date(2016, 9, 6, 0, 1, 0, 0, time.UTC)
	dec, err := currency.ConvertStringAt("1.0000", currency.USD, currency.EUR, at)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("1 USD to EUR = %s", dec.StringFixed(4))
	// Output:
	// 1 USD to EUR = 0.8961
}
