package bloomberg

import (
	"fmt"
	"time"

	"github.com/gopher-finance/forecast"
	"github.com/gopher-finance/forecast/models"
)

var VERBOSITY = 0

type BloombergSource struct {
}

func New() forecast.Source {
	return &BloombergSource{}
}

func (s *BloombergSource) Quote(symbols []string) (models.Quotes, error) {
	symbols = convertSymbols(symbols)

	quotes := make(models.Quotes)

	results := make(chan *models.Quote, len(symbols))
	errors := make(chan error, len(symbols))

	/* fetch all symbols in parallel */
	for _, symbol := range symbols {
		go func(symbol string) {
			quote, err := getQuote(symbol)
			if err != nil {
				errors <- err
			} else {
				results <- quote
			}
		}(symbol)
	}

	for i := 0; i < len(symbols); i++ {
		select {
		case err := <-errors:
			fmt.Println("bloomberg: error while fetching,", err)
		case r := <-results:
			r.Symbol = bloombergToYahoo(r.Symbol)
			quotes[r.Symbol] = r
		}
	}

	return quotes, nil
}

func (s *BloombergSource) Hist(symbols []string) (models.HistMap, error) {
	symbols = convertSymbols(symbols)

	m := make(models.HistMap, 0)

	results := make(chan *models.Hist, len(symbols))
	errors := make(chan error, len(symbols))

	/* fetch all symbols in parallel */
	for _, symbol := range symbols {
		go func(symbol string) {
			quote, err := getHist(symbol)
			if err != nil {
				errors <- err
			} else {
				results <- quote
			}
		}(symbol)
	}

	for i := 0; i < len(symbols); i++ {
		select {
		case err := <-errors:
			fmt.Println("bloomberg: hist error,", err)
		case r := <-results:
			r.Symbol = bloombergToYahoo(r.Symbol)
			m[r.Symbol] = r
		}
	}

	return m, nil
}

func (s *BloombergSource) HistLimit(symbols []string, start time.Time, end time.Time) (models.HistMap, error) {
	return nil, fmt.Errorf(ErrTplNotSupported, s.String(), "histlimit")
}

func (s *BloombergSource) DividendHist(symbols []string) (models.DividendHistMap, error) {
	return nil, fmt.Errorf(ErrTplNotSupported, s.String(), "dividendhist")
}

func (s *BloombergSource) DividendHistLimit(symbols []string, start time.Time, end time.Time) (models.DividendHistMap, error) {
	return nil, fmt.Errorf(ErrTplNotSupported, s.String(), "dividendhistlimi")
}

func (s *BloombergSource) String() string {
	return "Bloomberg"
}

func vprintln(a ...interface{}) (int, error) {
	if VERBOSITY > 0 {
		return fmt.Println(a...)
	}

	return 0, nil
}

func vprintf(format string, a ...interface{}) (int, error) {
	if VERBOSITY > 0 {
		return fmt.Printf(format, a...)
	}

	return 0, nil
}
