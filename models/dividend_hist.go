package models

type DividendHistMap map[string]*DividendHist

type DividendHist struct {
	Symbol    string
	Dividends []DividendEntry
}
