package other

import (
	"encoding/json"
	"errors"
	"fmt"
	"index-bot/logger"
	"index-bot/models"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MarketStackRest struct {
	AccessKey   string
	requestLeft uint16
	lockedTo    int64  // Timestamp (in ms) to when it's locked, 0 means there's no lock.
	fails       uint16 // If request failed, try again up to 3 times (delay 250/500/750ms) - after 3rd failed attempt => panic
}

// Collects all EOD stack entries from today back to Janurary 1st and returns average.
func (rest *MarketStackRest) RequestYTD(symbol string) (models.StockData, error) {
	year, month, day := time.Now().Date()
	baseURI := fmt.Sprintf(
		"/eod?symbols=%s&date_from=%d-01-01&date_to=%d-%s%d-%s%d&limit=1000",
		symbol,
		year,
		year,
		sif(month < 10, "0", ""),
		month,
		sif(day < 10, "0", ""),
		day,
	)
	var avgOpen, avgClose float32 = 0, 0
	var counter, offset uint = 0, 0

	for {
		raw, err := rest.Request(fmt.Sprintf("%s&offset=%d", baseURI, offset), nil)
		if err != nil {
			return models.StockData{}, err
		}

		res := models.HistoryStockData{}
		err = json.Unmarshal(raw, &res)
		if err != nil {
			return models.StockData{}, errors.New("failed to parse received data from MarketStack API")
		}

		if res.Pagination.Count == 0 || len(res.Data) == 0 {
			break
		}

		for _, stock := range res.Data {
			counter++
			avgOpen = avgOpen + (stock.Open-avgOpen)/float32(counter)
			avgClose = avgClose + (stock.Close-avgClose)/float32(counter)
		}

		offset += res.Pagination.Count
	}

	return models.StockData{
		Open:   avgOpen,
		Close:  avgClose,
		Symbol: symbol,
	}, nil
}

func (rest *MarketStackRest) RequestEOD(symbol string) (models.StockData, error) {
	raw, err := rest.Request("/tickers/"+symbol+"/eod/latest", nil)
	if err != nil {
		return models.StockData{}, err
	}

	res := models.StockData{}
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return models.StockData{}, errors.New("failed to parse received data from MarketStack API")
	}

	return res, nil
}

func (rest *MarketStackRest) Request(route string, jsonPayload interface{}) ([]byte, error) {
	now := time.Now().Unix()
	if rest.lockedTo != 0 && rest.lockedTo > now {
		offset := time.Second * time.Duration(rest.lockedTo-now)
		logger.Warn.Printf("Reached API rate limit! Waiting %s.", offset.String())
		time.Sleep(time.Second * time.Duration(rest.lockedTo-now))
	}

	format := "?"
	if strings.Contains(route, format) {
		format = "&"
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.marketstack.com/v1%s%saccess_key=%s", route, format, rest.AccessKey), nil)
	if err != nil {
		return nil, errors.New("failed to initialize new request: " + err.Error())
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		rest.fails++
		if rest.fails == 3 {
			logger.Error.Panicln("failed to make http request 3 times to https://api.marketstack.com/v1" + route + " (check internet connection and/or app credentials)")
		} else {
			time.Sleep(time.Millisecond * time.Duration(250*rest.fails))
			return rest.Request(route, jsonPayload) // Try again after potential internet connection failure.
		}
	}
	defer res.Body.Close()

	rest.fails = 0
	remaining, err := strconv.ParseFloat(res.Header.Get("x-quota-remaining"), 32)
	if err == nil && remaining == 0 {
		rest.lockedTo = now + 24*60
		rest.requestLeft = 0
	}
	rest.requestLeft = uint16(remaining)

	if res.StatusCode == 204 {
		return nil, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("failed to parse response body (json): " + err.Error())
	}

	if res.StatusCode >= 400 {
		return nil, errors.New(res.Status + " :: " + string(body))
	}

	return body, nil
}

func CreateMarketStackRest(accessKey string, limit uint16) MarketStackRest {
	return MarketStackRest{
		AccessKey:   accessKey,
		requestLeft: limit,
		lockedTo:    0,
		fails:       0,
	}
}
