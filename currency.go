// Copyright 2015 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// package currency implements a currency converter

package currency

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type Currency string

const (
	AED Currency = "AED" //	United Arab Emirates Dirham
	AFN          = "AFN" //	Afghanistan Afghani
	ALL          = "ALL" //	Albania Lek
	AMD          = "AMD" //	Armenia Dram
	ANG          = "ANG" //	Netherlands Antilles Guilder
	AOA          = "AOA" //	Angola Kwanza
	ARS          = "ARS" //	Argentina Peso
	AUD          = "AUD" //	Australia Dollar
	AWG          = "AWG" //	Aruba Guilder
	AZN          = "AZN" //	Azerbaijan New Manat
	BAM          = "BAM" //	Bosnia and Herzegovina Convertible Marka
	BBD          = "BBD" //	Barbados Dollar
	BDT          = "BDT" //	Bangladesh Taka
	BGN          = "BGN" //	Bulgaria Lev
	BHD          = "BHD" //	Bahrain Dinar
	BIF          = "BIF" //	Burundi Franc
	BMD          = "BMD" //	Bermuda Dollar
	BND          = "BND" //	Brunei Darussalam Dollar
	BOB          = "BOB" //	Bolivia Boliviano
	BRL          = "BRL" //	Brazil Real
	BSD          = "BSD" //	Bahamas Dollar
	BTN          = "BTN" //	Bhutan Ngultrum
	BWP          = "BWP" //	Botswana Pula
	BYR          = "BYR" //	Belarus Ruble
	BZD          = "BZD" //	Belize Dollar
	CAD          = "CAD" //	Canada Dollar
	CDF          = "CDF" //	Congo/Kinshasa Franc
	CHF          = "CHF" //	Switzerland Franc
	CLP          = "CLP" //	Chile Peso
	CNY          = "CNY" //	China Yuan Renminbi
	COP          = "COP" //	Colombia Peso
	CRC          = "CRC" //	Costa Rica Colon
	CUC          = "CUC" //	Cuba Convertible Peso
	CUP          = "CUP" //	Cuba Peso
	CVE          = "CVE" //	Cape Verde Escudo
	CZK          = "CZK" //	Czech Republic Koruna
	DJF          = "DJF" //	Djibouti Franc
	DKK          = "DKK" //	Denmark Krone
	DOP          = "DOP" //	Dominican Republic Peso
	DZD          = "DZD" //	Algeria Dinar
	EGP          = "EGP" //	Egypt Pound
	ERN          = "ERN" //	Eritrea Nakfa
	ETB          = "ETB" //	Ethiopia Birr
	EUR          = "EUR" //	Euro Member Countries
	FJD          = "FJD" //	Fiji Dollar
	FKP          = "FKP" //	Falkland Islands (Malvinas) Pound
	GBP          = "GBP" //	United Kingdom Pound
	GEL          = "GEL" //	Georgia Lari
	GGP          = "GGP" //	Guernsey Pound
	GHS          = "GHS" //	Ghana Cedi
	GIP          = "GIP" //	Gibraltar Pound
	GMD          = "GMD" //	Gambia Dalasi
	GNF          = "GNF" //	Guinea Franc
	GTQ          = "GTQ" //	Guatemala Quetzal
	GYD          = "GYD" //	Guyana Dollar
	HKD          = "HKD" //	Hong Kong Dollar
	HNL          = "HNL" //	Honduras Lempira
	HRK          = "HRK" //	Croatia Kuna
	HTG          = "HTG" //	Haiti Gourde
	HUF          = "HUF" //	Hungary Forint
	IDR          = "IDR" //	Indonesia Rupiah
	ILS          = "ILS" //	Israel Shekel
	IMP          = "IMP" //	Isle of Man Pound
	INR          = "INR" //	India Rupee
	IQD          = "IQD" //	Iraq Dinar
	IRR          = "IRR" //	Iran Rial
	ISK          = "ISK" //	Iceland Krona
	JEP          = "JEP" //	Jersey Pound
	JMD          = "JMD" //	Jamaica Dollar
	JOD          = "JOD" //	Jordan Dinar
	JPY          = "JPY" //	Japan Yen
	KES          = "KES" //	Kenya Shilling
	KGS          = "KGS" //	Kyrgyzstan Som
	KHR          = "KHR" //	Cambodia Riel
	KMF          = "KMF" //	Comoros Franc
	KPW          = "KPW" //	Korea (North) Won
	KRW          = "KRW" //	Korea (South) Won
	KWD          = "KWD" //	Kuwait Dinar
	KYD          = "KYD" //	Cayman Islands Dollar
	KZT          = "KZT" //	Kazakhstan Tenge
	LAK          = "LAK" //	Laos Kip
	LBP          = "LBP" //	Lebanon Pound
	LKR          = "LKR" //	Sri Lanka Rupee
	LRD          = "LRD" //	Liberia Dollar
	LSL          = "LSL" //	Lesotho Loti
	LYD          = "LYD" //	Libya Dinar
	MAD          = "MAD" //	Morocco Dirham
	MDL          = "MDL" //	Moldova Leu
	MGA          = "MGA" //	Madagascar Ariary
	MKD          = "MKD" //	Macedonia Denar
	MMK          = "MMK" //	Myanmar (Burma) Kyat
	MNT          = "MNT" //	Mongolia Tughrik
	MOP          = "MOP" //	Macau Pataca
	MRO          = "MRO" //	Mauritania Ouguiya
	MUR          = "MUR" //	Mauritius Rupee
	MVR          = "MVR" //	Maldives (Maldive Islands) Rufiyaa
	MWK          = "MWK" //	Malawi Kwacha
	MXN          = "MXN" //	Mexico Peso
	MYR          = "MYR" //	Malaysia Ringgit
	MZN          = "MZN" //	Mozambique Metical
	NAD          = "NAD" //	Namibia Dollar
	NGN          = "NGN" //	Nigeria Naira
	NIO          = "NIO" //	Nicaragua Cordoba
	NOK          = "NOK" //	Norway Krone
	NPR          = "NPR" //	Nepal Rupee
	NZD          = "NZD" //	New Zealand Dollar
	OMR          = "OMR" //	Oman Rial
	PAB          = "PAB" //	Panama Balboa
	PEN          = "PEN" //	Peru Nuevo Sol
	PGK          = "PGK" //	Papua New Guinea Kina
	PHP          = "PHP" //	Philippines Peso
	PKR          = "PKR" //	Pakistan Rupee
	PLN          = "PLN" //	Poland Zloty
	PYG          = "PYG" //	Paraguay Guarani
	QAR          = "QAR" //	Qatar Riyal
	RON          = "RON" //	Romania New Leu
	RSD          = "RSD" //	Serbia Dinar
	RUB          = "RUB" //	Russia Ruble
	RWF          = "RWF" //	Rwanda Franc
	SAR          = "SAR" //	Saudi Arabia Riyal
	SBD          = "SBD" //	Solomon Islands Dollar
	SCR          = "SCR" //	Seychelles Rupee
	SDG          = "SDG" //	Sudan Pound
	SEK          = "SEK" //	Sweden Krona
	SGD          = "SGD" //	Singapore Dollar
	SHP          = "SHP" //	Saint Helena Pound
	SLL          = "SLL" //	Sierra Leone Leone
	SOS          = "SOS" //	Somalia Shilling
	SPL          = "SPL" // Seborga Luigino
	SRD          = "SRD" //	Suriname Dollar
	STD          = "STD" //	São Tomé and Príncipe Dobra
	SVC          = "SVC" //	El Salvador Colon
	SYP          = "SYP" //	Syria Pound
	SZL          = "SZL" //	Swaziland Lilangeni
	THB          = "THB" //	Thailand Baht
	TJS          = "TJS" //	Tajikistan Somoni
	TMT          = "TMT" //	Turkmenistan Manat
	TND          = "TND" //	Tunisia Dinar
	TOP          = "TOP" //	Tonga Pa'anga
	TRY          = "TRY" //	Turkey Lira
	TTD          = "TTD" //	Trinidad and Tobago Dollar
	TVD          = "TVD" //	Tuvalu Dollar
	TWD          = "TWD" //	Taiwan New Dollar
	TZS          = "TZS" //	Tanzania Shilling
	UAH          = "UAH" //	Ukraine Hryvnia
	UGX          = "UGX" //	Uganda Shilling
	USD          = "USD" //	United States Dollar
	UYU          = "UYU" //	Uruguay Peso
	UZS          = "UZS" //	Uzbekistan Som
	VEF          = "VEF" //	Venezuela Bolivar
	VND          = "VND" //	Viet Nam Dong
	VUV          = "VUV" //	Vanuatu Vatu
	WST          = "WST" //	Samoa Tala
	XAF          = "XAF" //	Communauté Financière Africaine (BEAC) CFA Franc BEAC
	XCD          = "XCD" //	East Caribbean Dollar
	XDR          = "XDR" //	International Monetary Fund (IMF) Special Drawing Rights
	XOF          = "XOF" //	Communauté Financière Africaine (BCEAO) Franc
	XPF          = "XPF" //	Comptoirs Français du Pacifique (CFP) Franc
	YER          = "YER" //	Yemen Rial
	ZAR          = "ZAR" //	South Africa Rand
	ZMW          = "ZMW" //	Zambia Kwacha
	ZWD          = "ZWD" //	Zimbabwe Dollar
)

var currencies = [...]Currency{
	AED, AFN, ALL, AMD, ANG, AOA, ARS, AUD, AWG, AZN, BAM, BBD, BDT, BGN, BHD,
	BIF, BMD, BND, BOB, BRL, BSD, BTN, BWP, BYR, BZD, CAD, CDF, CHF, CLP, CNY,
	COP, CRC, CUC, CUP, CVE, CZK, DJF, DKK, DOP, DZD, EGP, ERN, ETB, EUR, FJD,
	FKP, GBP, GEL, GGP, GHS, GIP, GMD, GNF, GTQ, GYD, HKD, HNL, HRK, HTG, HUF,
	IDR, ILS, IMP, INR, IQD, IRR, ISK, JEP, JMD, JOD, JPY, KES, KGS, KHR, KMF,
	KPW, KRW, KWD, KYD, KZT, LAK, LBP, LKR, LRD, LSL, LYD, MAD, MDL, MGA, MKD,
	MMK, MNT, MOP, MRO, MUR, MVR, MWK, MXN, MYR, MZN, NAD, NGN, NIO, NOK, NPR,
	NZD, OMR, PAB, PEN, PGK, PHP, PKR, PLN, PYG, QAR, RON, RSD, RUB, RWF, SAR,
	SBD, SCR, SDG, SEK, SGD, SHP, SLL, SOS, SPL, SRD, STD, SVC, SYP, SZL, THB,
	TJS, TMT, TND, TOP, TRY, TTD, TVD, TWD, TZS, UAH, UGX, USD, UYU, UZS, VEF,
	VND, VUV, WST, XAF, XCD, XDR, XOF, XPF, YER, ZAR, ZMW, ZWD}

type Converter struct {
	ex  map[Currency]decimal.Decimal
	mux sync.Mutex
}

func New() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertString(value string, from, to Currency) decimal.Decimal {
	v, _ := decimal.NewFromString(value)
	return c.genConvert(v, from, to, nil)
}

func (c *Converter) ConvertStringAt(value string, from, to Currency, at time.Time) decimal.Decimal {
	v, _ := decimal.NewFromString(value)
	return c.genConvert(v, from, to, &at)
}

func (c *Converter) Convert(value decimal.Decimal, from, to Currency) decimal.Decimal {
	return c.genConvert(value, from, to, nil)
}

func (c *Converter) ConvertAt(value decimal.Decimal, from, to Currency, at time.Time) decimal.Decimal {
	return c.genConvert(value, from, to, &at)
}

func (c *Converter) genConvert(value decimal.Decimal, from, to Currency, at *time.Time) decimal.Decimal {
	if to != EUR {
		// we only deal in EUR atm
		return decimal.Zero
	}

	c.mux.Lock()
	rate, ok := c.ex[from]
	c.mux.Unlock()

	if !ok {
		return decimal.Zero
	}

	//fmt.Println(value, from, "*", rate, to)
	return value.Mul(rate)
}

var oneD = decimal.NewFromFloat(1.0)

func (c *Converter) Update() error {
	// create URL
	var urlbuffer bytes.Buffer

	// "http://download.finance.yahoo.com/d/quotes.csv?f=snl1d1t1ab&s= + &EURPLN=X
	urlbuffer.Grow(62 + (len(currencies) * 9) - 1)

	urlbuffer.WriteString("http://download.finance.yahoo.com/d/quotes.csv?f=snl1d1t1ab&s=")
	urlbuffer.WriteString(string(currencies[0]) + "EUR" + "=X")

	for _, s := range currencies[1:] {
		urlbuffer.WriteString("," + string(s) + "EUR=X")
	}

	// do request
	res, err := readCSVFromUrl(urlbuffer.String())

	if err != nil {
		return err
	}

	// parse response
	ex := make(map[Currency]decimal.Decimal, len(res))

	for _, row := range res {
		v := row[2]

		if _, err := strconv.ParseFloat(v, 64); err != nil {
			continue
		}

		rate, err := decimal.NewFromString(v)

		if err != nil {
			// log
			continue
		}

		v = row[0]

		if len(v) != 8 {
			continue
		}

		cur := Currency(v[0:3])
		ex[cur] = rate
	}

	// update map
	c.mux.Lock()
	c.ex = ex
	c.mux.Unlock()
	return nil
}

func readCSVFromUrl(url string) ([][]string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	data, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	return data, nil
}

var DefaultConverter = New()

func ConvertString(value string, from, to Currency) decimal.Decimal {
	return DefaultConverter.ConvertString(value, from, to)
}

func ConvertStringAt(value string, from, to Currency, at time.Time) decimal.Decimal {
	return DefaultConverter.ConvertStringAt(value, from, to, at)
}

func Convert(value decimal.Decimal, from, to Currency) decimal.Decimal {
	return DefaultConverter.Convert(value, from, to)
}

func ConvertAt(value decimal.Decimal, from, to Currency, at time.Time) decimal.Decimal {
	return DefaultConverter.ConvertAt(value, from, to, at)
}

//http://download.finance.yahoo.com/d/quotes.csv?s=EURPLN=X&f=snl1d1t1ab
