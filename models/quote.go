package models

import (
	"encoding/json"
	"math"
	"time"
)

const (
	ErrTplNotSupported = "source '%s' does not support action '%s'"
)

type Quotes map[string]*Quote

type Quote struct {
	Symbol   string /* e.g.: VEUR.AS, Vanguard dev. europe on Amsterdam */
	Name     string
	Exchange string

	/* last actualization of the results */
	Updated time.Time

	/* volume */
	Volume         int64 /* outstanding shares */
	AvgDailyVolume int64 /* avg amount of shares traded */

	/* dividend & related */
	PeRatio          float64   /* Price / EPS */
	EarningsPerShare float64   /* (net income - spec.dividends) / avg.  outstanding shares */
	DividendPerShare float64   /* total (non-special) dividend payout / total amount of shares */
	DividendYield    float64   /* annual div. per share / price per share */
	DividendExDate   time.Time /* last dividend payout date */

	/* price & derived */
	Bid, Ask              float64
	Open, PreviousClose   float64
	LastTradePrice        float64
	Change, ChangePercent float64

	DayLow, DayHigh   float64
	YearLow, YearHigh float64

	Ma50, Ma200 float64 /* 200- and 50-day moving average */
}

/* will try to calculate the dividend payout ratio, if possible,
 * otherwise returns 0 */
func (q *Quote) DivPayoutRatio() float64 {
	/* total dividends / net income (same period, mostly 1Y):
	 * TODO: implement this (first implement historical data
	 * aggregation) */

	/* DPS / EPS */
	if q.DividendPerShare != 0 && q.EarningsPerShare != 0 {
		return q.DividendPerShare / q.EarningsPerShare
	}
	return 0
}

func (q *Quote) GetPrice() float64 {
	if q.Ask != 0 {
		return q.Ask
	}
	return q.LastTradePrice
}

// SharesToBuy gives you the number of shares to buy if you want
//  * the transaction cost to be less than a certain percentage
func (q *Quote) SharesToBuy(txCost, desiredTxCostPerc float64) float64 {
	return math.Ceil((txCost - desiredTxCostPerc*txCost) /
		(desiredTxCostPerc * q.GetPrice()))
}

// True if the Quote is increasing
// False otherwise
func (q *Quote) IsIncreasing() bool {
	return q.LastTradePrice >= q.PreviousClose
}

func (q *Quote) VariationValue() float64 {
	return q.LastTradePrice - q.PreviousClose
}

func (q *Quote) VariationValuePercent() float64 {
	return q.VariationValue() / q.PreviousClose * 100
}

func (q *Quote) WouldBuyOrSell() interface{} {
	return q.PreviousClose > q.Ma200
}

func (q *Quote) JSON() string {
	b, _ := json.Marshal(q)
	return string(b)
}
