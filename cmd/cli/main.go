package main

import (
	"fmt"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/gopher-finance/forecast"
	"github.com/gopher-finance/forecast/sources/bloomberg"
	"github.com/gopher-finance/forecast/sources/yahoofinance"
)

// Variables linked during build
var (
	Version    string
	Githash    string
	Buildstamp string
	AppName    string
	SourcesMap map[string]forecast.Source

	symbols = kingpin.Flag("symbols", "ticker symbols to analyze").Short('s').Required().Strings()
	source  = kingpin.Flag("source", "bloomberg || yahoo_yql || yahoo_csv").
		Default("yahoo_yql").Enum("bloomberg", "yahoo_yql", "yahoo_csv")
)

func init() {
	SourcesMap = map[string]forecast.Source{
		"bloomberg": bloomberg.New(),
		"yahoo_yql": yahoofinance.NewYql(),
		"yahoo_csv": yahoofinance.NewCsv()}
}

func main() {
	kingpin.Version(fmt.Sprintf("%s [%s] Commit: %s Buildtstamp: %s",
		AppName, Version, Githash, Buildstamp))
	kingpin.Parse()

	app := App{Symbols: *symbols, Src: SourcesMap[*source]}
	app.Run()
}
