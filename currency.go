// Copyright 2015 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// package currency implements a currency converter
package currency

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

// Currency is a ISO 4217 three-letter alphabetic code.
type Currency string

// Scan implements the Scanner interface.
// The value type must be string / []byte otherwise Scan fails.
func (c *Currency) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("nil err")
	}

	var s string

	switch v := value.(type) {
	case []byte:
		s = string(v)
	case string:
		s = v
	default:
		return fmt.Errorf("Can't convert %T to currency.Currency", value)
	}

	cc, err := ParseCurrency(s)

	if err != nil {
		return err
	}

	*c = cc
	return nil
}

// Value implements the driver Valuer interface.
func (c Currency) Value() (driver.Value, error) {
	return string(c), nil
}

const (
	AED Currency = "AED" // United Arab Emirates Dirham
	AFN          = "AFN" // Afghanistan Afghani
	ALL          = "ALL" // Albania Lek
	AMD          = "AMD" // Armenia Dram
	ANG          = "ANG" // Netherlands Antilles Guilder
	AOA          = "AOA" // Angola Kwanza
	ARS          = "ARS" // Argentina Peso
	AUD          = "AUD" // Australia Dollar
	AWG          = "AWG" // Aruba Guilder
	AZN          = "AZN" // Azerbaijan New Manat
	BAM          = "BAM" // Bosnia and Herzegovina Convertible Marka
	BBD          = "BBD" // Barbados Dollar
	BDT          = "BDT" // Bangladesh Taka
	BGN          = "BGN" // Bulgaria Lev
	BHD          = "BHD" // Bahrain Dinar
	BIF          = "BIF" // Burundi Franc
	BMD          = "BMD" // Bermuda Dollar
	BND          = "BND" // Brunei Darussalam Dollar
	BOB          = "BOB" // Bolivia Boliviano
	BRL          = "BRL" // Brazil Real
	BSD          = "BSD" // Bahamas Dollar
	BTN          = "BTN" // Bhutan Ngultrum
	BWP          = "BWP" // Botswana Pula
	BYR          = "BYR" // Belarus Ruble
	BZD          = "BZD" // Belize Dollar
	CAD          = "CAD" // Canada Dollar
	CDF          = "CDF" // Congo/Kinshasa Franc
	CHF          = "CHF" // Switzerland Franc
	CLP          = "CLP" // Chile Peso
	CNY          = "CNY" // China Yuan Renminbi
	COP          = "COP" // Colombia Peso
	CRC          = "CRC" // Costa Rica Colon
	CUC          = "CUC" // Cuba Convertible Peso
	CUP          = "CUP" // Cuba Peso
	CVE          = "CVE" // Cape Verde Escudo
	CZK          = "CZK" // Czech Republic Koruna
	DJF          = "DJF" // Djibouti Franc
	DKK          = "DKK" // Denmark Krone
	DOP          = "DOP" // Dominican Republic Peso
	DZD          = "DZD" // Algeria Dinar
	EGP          = "EGP" // Egypt Pound
	ERN          = "ERN" // Eritrea Nakfa
	ETB          = "ETB" // Ethiopia Birr
	EUR          = "EUR" // Euro Member Countries
	FJD          = "FJD" // Fiji Dollar
	FKP          = "FKP" // Falkland Islands (Malvinas) Pound
	GBP          = "GBP" // United Kingdom Pound
	GEL          = "GEL" // Georgia Lari
	GGP          = "GGP" // Guernsey Pound
	GHS          = "GHS" // Ghana Cedi
	GIP          = "GIP" // Gibraltar Pound
	GMD          = "GMD" // Gambia Dalasi
	GNF          = "GNF" // Guinea Franc
	GTQ          = "GTQ" // Guatemala Quetzal
	GYD          = "GYD" // Guyana Dollar
	HKD          = "HKD" // Hong Kong Dollar
	HNL          = "HNL" // Honduras Lempira
	HRK          = "HRK" // Croatia Kuna
	HTG          = "HTG" // Haiti Gourde
	HUF          = "HUF" // Hungary Forint
	IDR          = "IDR" // Indonesia Rupiah
	ILS          = "ILS" // Israel Shekel
	IMP          = "IMP" // Isle of Man Pound
	INR          = "INR" // India Rupee
	IQD          = "IQD" // Iraq Dinar
	IRR          = "IRR" // Iran Rial
	ISK          = "ISK" // Iceland Krona
	JEP          = "JEP" // Jersey Pound
	JMD          = "JMD" // Jamaica Dollar
	JOD          = "JOD" // Jordan Dinar
	JPY          = "JPY" // Japan Yen
	KES          = "KES" // Kenya Shilling
	KGS          = "KGS" // Kyrgyzstan Som
	KHR          = "KHR" // Cambodia Riel
	KMF          = "KMF" // Comoros Franc
	KPW          = "KPW" // Korea (North) Won
	KRW          = "KRW" // Korea (South) Won
	KWD          = "KWD" // Kuwait Dinar
	KYD          = "KYD" // Cayman Islands Dollar
	KZT          = "KZT" // Kazakhstan Tenge
	LAK          = "LAK" // Laos Kip
	LBP          = "LBP" // Lebanon Pound
	LKR          = "LKR" // Sri Lanka Rupee
	LRD          = "LRD" // Liberia Dollar
	LSL          = "LSL" // Lesotho Loti
	LYD          = "LYD" // Libya Dinar
	MAD          = "MAD" // Morocco Dirham
	MDL          = "MDL" // Moldova Leu
	MGA          = "MGA" // Madagascar Ariary
	MKD          = "MKD" // Macedonia Denar
	MMK          = "MMK" // Myanmar (Burma) Kyat
	MNT          = "MNT" // Mongolia Tughrik
	MOP          = "MOP" // Macau Pataca
	MRO          = "MRO" // Mauritania Ouguiya
	MUR          = "MUR" // Mauritius Rupee
	MVR          = "MVR" // Maldives (Maldive Islands) Rufiyaa
	MWK          = "MWK" // Malawi Kwacha
	MXN          = "MXN" // Mexico Peso
	MXV          = "MXV" // Mexican Unidad de Inversion (UDI) (funds code)
	MYR          = "MYR" // Malaysia Ringgit
	MZN          = "MZN" // Mozambique Metical
	NAD          = "NAD" // Namibia Dollar
	NGN          = "NGN" // Nigeria Naira
	NIO          = "NIO" // Nicaragua Cordoba
	NOK          = "NOK" // Norway Krone
	NPR          = "NPR" // Nepal Rupee
	NZD          = "NZD" // New Zealand Dollar
	OMR          = "OMR" // Oman Rial
	PAB          = "PAB" // Panama Balboa
	PEN          = "PEN" // Peru Nuevo Sol
	PGK          = "PGK" // Papua New Guinea Kina
	PHP          = "PHP" // Philippines Peso
	PKR          = "PKR" // Pakistan Rupee
	PLN          = "PLN" // Poland Zloty
	PYG          = "PYG" // Paraguay Guarani
	QAR          = "QAR" // Qatar Riyal
	RON          = "RON" // Romania New Leu
	RSD          = "RSD" // Serbia Dinar
	RUB          = "RUB" // Russia Ruble
	RWF          = "RWF" // Rwanda Franc
	SAR          = "SAR" // Saudi Arabia Riyal
	SBD          = "SBD" // Solomon Islands Dollar
	SCR          = "SCR" // Seychelles Rupee
	SDG          = "SDG" // Sudan Pound
	SEK          = "SEK" // Sweden Krona
	SGD          = "SGD" // Singapore Dollar
	SHP          = "SHP" // Saint Helena Pound
	SLL          = "SLL" // Sierra Leone Leone
	SOS          = "SOS" // Somalia Shilling
	SPL          = "SPL" // Seborga Luigino
	SRD          = "SRD" // Suriname Dollar
	STD          = "STD" // São Tomé and Príncipe Dobra
	SVC          = "SVC" // El Salvador Colon
	SYP          = "SYP" // Syria Pound
	SZL          = "SZL" // Swaziland Lilangeni
	THB          = "THB" // Thailand Baht
	TJS          = "TJS" // Tajikistan Somoni
	TMT          = "TMT" // Turkmenistan Manat
	TND          = "TND" // Tunisia Dinar
	TOP          = "TOP" // Tonga Pa'anga
	TRY          = "TRY" // Turkey Lira
	TTD          = "TTD" // Trinidad and Tobago Dollar
	TVD          = "TVD" // Tuvalu Dollar
	TWD          = "TWD" // Taiwan New Dollar
	TZS          = "TZS" // Tanzania Shilling
	UAH          = "UAH" // Ukraine Hryvnia
	UGX          = "UGX" // Uganda Shilling
	USD          = "USD" // United States Dollar
	UYU          = "UYU" // Uruguay Peso
	UZS          = "UZS" // Uzbekistan Som
	VEF          = "VEF" // Venezuela Bolivar
	VND          = "VND" // Viet Nam Dong
	VUV          = "VUV" // Vanuatu Vatu
	WST          = "WST" // Samoa Tala
	XAF          = "XAF" // Communauté Financière Africaine (BEAC) CFA Franc BEAC
	XCD          = "XCD" // East Caribbean Dollar
	XDR          = "XDR" // International Monetary Fund (IMF) Special Drawing Rights
	XOF          = "XOF" // Communauté Financière Africaine (BCEAO) Franc
	XPF          = "XPF" // Comptoirs Français du Pacifique (CFP) Franc
	YER          = "YER" // Yemen Rial
	ZAR          = "ZAR" // South Africa Rand
	ZMW          = "ZMW" // Zambia Kwacha
	ZWD          = "ZWD" // Zimbabwe Dollar

	// Metals
	XAG Currency = "XAU" // Gold
	XAU          = "XAG" // Silver
	XCP          = "XCP" // Copper
	XPD          = "XPD" // Palladium
	XPT          = "XPT" // Platinum

	// Historic currencies
	CYP Currency = "CYP" // Cypriot pound
	DEM          = "DEM" // German mark
	ECS          = "ECS" // Ecuadorian sucre
	FRF          = "FRF" // French franc
	IEP          = "IEP" // Irish pound (punt in Irish language)
	ITL          = "ITL" // Italian lira
	LTL          = "LTL" // Lithuanian litas
	LVL          = "LVL" // Latvian lats
	SIT          = "SIT" // Slovenian tolar
	ZWL          = "ZWL" // Zimbabwean dollar A/10

	// Unofficial currency codes
	CNH Currency = "CNH" // Chinese yuan (when traded offshore)

	// Other
	CLF Currency = "CLF" // Unidad de Fomento (funds code)

)

var currencies = [...]Currency{
	AED, AFN, ALL, AMD, ANG, AOA, ARS, AUD, AWG, AZN, BAM, BBD, BDT, BGN, BHD,
	BIF, BMD, BND, BOB, BRL, BSD, BTN, BWP, BYR, BZD, CAD, CDF, CHF, CLP, CNY,
	COP, CRC, CUC, CUP, CVE, CZK, DJF, DKK, DOP, DZD, EGP, ERN, ETB, EUR, FJD,
	FKP, GBP, GEL, GGP, GHS, GIP, GMD, GNF, GTQ, GYD, HKD, HNL, HRK, HTG, HUF,
	IDR, ILS, IMP, INR, IQD, IRR, ISK, JEP, JMD, JOD, JPY, KES, KGS, KHR, KMF,
	KPW, KRW, KWD, KYD, KZT, LAK, LBP, LKR, LRD, LSL, LYD, MAD, MDL, MGA, MKD,
	MMK, MNT, MOP, MRO, MUR, MVR, MWK, MXN, MXV, MYR, MZN, NAD, NGN, NIO, NOK,
	NPR, NZD, OMR, PAB, PEN, PGK, PHP, PKR, PLN, PYG, QAR, RON, RSD, RUB, RWF,
	SAR, SBD, SCR, SDG, SEK, SGD, SHP, SLL, SOS, SPL, SRD, STD, SVC, SYP, SZL,
	THB, TJS, TMT, TND, TOP, TRY, TTD, TVD, TWD, TZS, UAH, UGX, USD, UYU, UZS,
	VEF, VND, VUV, WST, XAF, XCD, XDR, XOF, XPF, YER, ZAR, ZMW, ZWD, XAU, XAG,
	XCP, XPD, XPT, CYP, DEM, ECS, FRF, IEP, ITL, LTL, LVL, SIT, ZWL, CNH, CLF}

var ErrCurrencyLength = errors.New("Currency should be 3 char long")
var ErrCurrencyUnknown = errors.New("Currency is unknown")
var ErrNotExist = errors.New("Exchange rate or Currency does not exist")
var ErrFetchingData = errors.New("Unable to fetch data for date")

// ParseCurrency returns the Currency value represented by the string.
func ParseCurrency(v string) (Currency, error) {
	if len(v) != 3 {
		return "", ErrCurrencyLength
	}

	cur := Currency(strings.ToUpper(v))

	for _, curb := range currencies {
		if curb == cur {
			return cur, nil
		}
	}

	return "", ErrCurrencyUnknown
}

// Converter holds an Exchange and implements some methods for currency
// conversion using the exchange.
type Converter struct {
	ex *Exchange
}

// New initializes an Converter
func New() *Converter {
	return &Converter{
		ex: NewExchange(),
	}
}

// Convert converts the decimal value to the given currency.
func (c *Converter) Convert(value decimal.Decimal, from, to Currency) (decimal.Decimal, error) {
	return c.genConvert(value, from, to, nil)
}

// ConvertAt converts the decimal value to the given currency using the
// exchange rate from the date specificied.
func (c *Converter) ConvertAt(value decimal.Decimal, from, to Currency, at time.Time) (decimal.Decimal, error) {
	return c.genConvert(value, from, to, &at)
}

// ConvertString converts the decimal value (represented as a string) to the
// given currency.
func (c *Converter) ConvertString(value string, from, to Currency) (decimal.Decimal, error) {
	v, _ := decimal.NewFromString(value)
	return c.genConvert(v, from, to, nil)
}

// ConvertStringAt converts the decimal value (represented as a string) to the
// given currency using the exchange rate from the date specificied.
func (c *Converter) ConvertStringAt(value string, from, to Currency, at time.Time) (decimal.Decimal, error) {
	v, _ := decimal.NewFromString(value)
	return c.genConvert(v, from, to, &at)
}

func (c *Converter) genConvert(value decimal.Decimal, from, to Currency, at *time.Time) (decimal.Decimal, error) {
	var t time.Time

	if at == nil {
		t = time.Now().UTC()
	} else {
		t = *at
	}

	var usd decimal.Decimal

	if from == USD {
		usd = value
	} else {
		fromRate, err := c.ex.Get(t, from)

		if err != nil {
			return decimal.Zero, err
		}

		usd = value.Mul(fromRate.ToUSD)
	}

	if to == USD {
		return usd, nil
	}

	toRate, err := c.ex.Get(t, to)

	if err != nil {
		return decimal.Zero, err
	}

	return usd.Mul(toRate.FromUSD), nil
}

// DefaultConverter is the default Converter and is used by Convert, ConvertAt,
// ConvertString and ConvertStringAt.
var DefaultConverter = New()

// Convert converts the decimal value to the given currency.
//
// Convert is a wrapper for DefaultConverter.Convert.
func Convert(value decimal.Decimal, from, to Currency) (decimal.Decimal, error) {
	return DefaultConverter.Convert(value, from, to)
}

// ConvertAt converts the decimal value to the given currency using the
// exchange rate from the date specificied.
//
// ConvertAt is a wrapper for DefaultConverter.ConvertAt.
func ConvertAt(value decimal.Decimal, from, to Currency, at time.Time) (decimal.Decimal, error) {
	return DefaultConverter.ConvertAt(value, from, to, at)
}

// ConvertString converts the decimal value (represented as a string) to the
// given currency.
//
// ConvertString is a wrapper for DefaultConverter.ConvertString.
func ConvertString(value string, from, to Currency) (decimal.Decimal, error) {
	return DefaultConverter.ConvertString(value, from, to)
}

// ConvertStringAt converts the decimal value (represented as a string) to the
// given currency using the exchange rate from the date specificied.
//
// ConvertStringAt is a wrapper for DefaultConverter.ConvertStringAt.
func ConvertStringAt(value string, from, to Currency, at time.Time) (decimal.Decimal, error) {
	return DefaultConverter.ConvertStringAt(value, from, to, at)
}
