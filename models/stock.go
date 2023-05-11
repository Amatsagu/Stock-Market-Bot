package models

import (
	"image/color"
	"strconv"
)

type StockData struct {
	Open float32 `json:"open"`
	// High   float32 `json:"high"`
	// Low    float32 `json:"low"`
	Close float32 `json:"close"`
	// Volume float32 `json:"volume"`
	// Date   string  `json:"date"`
	Symbol string `json:"symbol"`
}

// Returns difference between open & close values - whether it's worth more or less (in %).
func (stock StockData) Diff() float64 {
	return float64(stock.Close)/float64(stock.Open)*100 - 100
}

func (stock StockData) DiffString() string {
	v := stock.Diff()
	if v >= 0 {
		return "+" + strconv.FormatFloat(v, 'f', 2, 64) + "%"
	} else {
		return strconv.FormatFloat(v, 'f', 2, 64) + "%"
	}
}

func (stock StockData) DiffPaint() (string, color.RGBA) {
	v := stock.Diff()
	if v >= 0 {
		return "+" + strconv.FormatFloat(v, 'f', 2, 64) + "%", color.RGBA{0, 111, 0, 200}
	} else {
		return strconv.FormatFloat(v, 'f', 2, 64) + "%", color.RGBA{168, 0, 0, 200}
	}
}

type HistoryStockData struct {
	Pagination Pagination
	Data       []StockData
}

type Pagination struct {
	Limit  uint `json:"limit"`
	Offset uint `json:"offset"`
	Count  uint `json:"count"`
	Total  uint `json:"total"`
}
