package models

import "github.com/gopher-finance/forecaster/utils"

type DividendEntry struct {
	Date      utils.YearMonthDay
	Dividends float64 `json:",string"`
}
