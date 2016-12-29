package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/gopher-finance/forecast"
)

type App struct {
	Symbols []string
	Src     forecast.Source
}

func (a *App) Run() {
	var buffer bytes.Buffer
	buffer.WriteString("requesting information on ")
	buffer.WriteString(strings.Join(a.Symbols, " "))
	buffer.WriteString("\n")
	res, err := a.Src.Quote(a.Symbols)
	if err != nil {
		log.Fatal("gofinance: could not fetch, ", err)
	}

	desiredTxCostPerc := 0.01
	txCost := 9.75

	for _, r := range res {
		buffer.WriteString(DisplayQuote(r, txCost, desiredTxCostPerc))
	}
	fmt.Println(buffer.String())
}
