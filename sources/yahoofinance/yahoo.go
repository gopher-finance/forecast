package yahoofinance

import (
	"fmt"
	"time"

	"github.com/gopher-finance/forecaster"
	"github.com/gopher-finance/forecaster/models"
)

const (
	HistoricalUrl = "http://ichart.finance.yahoo.com/table.csv"
)

const (
	TypeCsv = iota
	TypeYql
)

type Source struct {
	srcType int
}

func NewCsv() forecaster.Source {
	return &Source{srcType: TypeCsv}
}

func NewYql() forecaster.Source {
	return &Source{srcType: TypeYql}
}

func (s *Source) Quote(symbols []string) (models.Quotes, error) {
	switch s.srcType {
	case TypeCsv:
		return csvQuotes(symbols)
	case TypeYql:
		return yqlQuotes(symbols)
	}

	return nil, fmt.Errorf("yahoo finance: unknown backend type: %v", s.srcType)
}

func (s *Source) Hist(symbols []string) (models.HistMap, error) {
	return yqlHist(symbols, nil, nil)
}

func (s *Source) HistLimit(symbols []string, start time.Time, end time.Time) (models.HistMap, error) {
	return yqlHist(symbols, &start, &end)
}

func (s *Source) DividendHist(symbols []string) (models.DividendHistMap, error) {
	return yqlDivHist(symbols, nil, nil)
}

func (s *Source) DividendHistLimit(symbols []string, start time.Time, end time.Time) (models.DividendHistMap, error) {
	return yqlDivHist(symbols, &start, &end)
}

func (s *Source) String() string {
	return "Yahoo Finance"
}
