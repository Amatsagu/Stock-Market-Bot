package models

import (
	"image/color"
	"strconv"
)

type HistoricalMarketStockData struct {
	Symbol     string `json:"symbol"`
	Historical []StockData
}

type StockData struct {
	Date           string  `json:"date"`
	Open           float32 `json:"open"`
	High           float32 `json:"high"`
	Low            float32 `json:"low"`
	Close          float32 `json:"close"`
	Change         float32 `json:"change"`        // close / open
	ChangePercent  float32 `json:"changePercent"` // close / open * 100 - 100
	ChangeOverTime float32 `json:"changeOverTime"`
}

func (data HistoricalMarketStockData) DiffEODPaint() (string, color.RGBA) {
	v := float64(data.Historical[0].Close/data.Historical[1].Close*100 - 100)
	if v >= 0 {
		return "+" + strconv.FormatFloat(v, 'f', 2, 32) + "%", color.RGBA{0, 111, 0, 200}
	} else {
		return strconv.FormatFloat(v, 'f', 2, 32) + "%", color.RGBA{168, 0, 0, 200}
	}
}

// Returns diff from entire slice (YEAR.01.01 open -> YEAR.MM.DD close)
func (data HistoricalMarketStockData) DiffYTDPaint() (string, color.RGBA) {
	v := float64(data.Historical[0].Open/data.Historical[len(data.Historical)-1].Close*100 - 100)
	if v >= 0 {
		return "+" + strconv.FormatFloat(v, 'f', 2, 32) + "%", color.RGBA{0, 111, 0, 200}
	} else {
		return strconv.FormatFloat(v, 'f', 2, 32) + "%", color.RGBA{168, 0, 0, 200}
	}
}
