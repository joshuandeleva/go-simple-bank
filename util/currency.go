package util

// declare the supported currencies

const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
	JPY = "JPY"
	CNY = "CNY"
	KSH = "KSH"
	NGN = "NGN"
	CAD = "CAD"
	AUD = "AUD"
	INR = "INR"
	ZAR = "ZAR"
	CHF = "CHF"
	SEK = "SEK"
	NOK = "NOK"
	DKK = "DKK"
	TRY = "TRY"
	RUB = "RUB"
	BRL = "BRL"
	MXN = "MXN"
	ARS = "ARS"
	PLN = "PLN"
	CLP = "CLP"
	COP = "COP"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, GBP, JPY, CNY, KSH, NGN, CAD, AUD, INR, ZAR, CHF, SEK, NOK, DKK, TRY, RUB, BRL, MXN, ARS, PLN, CLP, COP:
		return true
	default:
		return false
	}
}
