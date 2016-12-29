package models

import (
	"encoding/json"
	"time"
)

type HistMap map[string]*Hist

type Hist struct {
	Symbol  string
	From    time.Time
	To      time.Time
	Entries []*HistEntry
}

func (h *Hist) MovingAverage() float64 {
	if len(h.Entries) == 0 {
		return 0
	}

	var sum, count float64

	for _, row := range h.Entries {
		sum += row.Close
		count++
	}

	return sum / count
}

func (h *Hist) JSON() string {
	b, _ := json.Marshal(h)
	return string(b)
}
