package other

import (
	"encoding/json"
	"errors"
	"fmt"
	"index-bot/logger"
	"index-bot/models"
	"io"
	"net/http"
	"sync"
	"time"
)

type MarketStackRest struct {
	mu          sync.RWMutex
	AccessKey   string
	requestLeft uint16
	lockedTo    int64  // Timestamp (in ms) to when it's locked, 0 means there's no lock.
	fails       uint16 // If request failed, try again up to 3 times (delay 250/500/750ms) - after 3rd failed attempt => panic
}

func (rest *MarketStackRest) RequestYearly(symbol string) (models.HistoricalMarketStockData, error) {
	year, month, day := time.Now().Date()
	today := models.DateParam{
		Year:  year,
		Month: month,
		Day:   day,
	}
	past := models.DateParam{
		Year:  year,
		Month: time.January,
		Day:   1,
	}

	raw, err := rest.Request(symbol, past, today)
	if err != nil {
		return models.HistoricalMarketStockData{}, err
	}

	res := models.HistoricalMarketStockData{}
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return models.HistoricalMarketStockData{}, errors.New("failed to parse received data from MarketStack API")
	}

	return res, nil
}

func (rest *MarketStackRest) RequestEOD(symbol string) (models.HistoricalMarketStockData, error) {
	year, month, day := time.Now().Date()
	today := models.DateParam{
		Year:  year,
		Month: month,
		Day:   day,
	}

	year2, month2, day2 := time.Now().Add(-(time.Hour * 24)).Date()
	past := models.DateParam{
		Year:  year2,
		Month: month2,
		Day:   day2,
	}

	raw, err := rest.Request(symbol, past, today)
	if err != nil {
		return models.HistoricalMarketStockData{}, err
	}

	res := models.HistoricalMarketStockData{}
	err = json.Unmarshal(raw, &res)
	if err != nil {
		return models.HistoricalMarketStockData{}, errors.New("failed to parse received data from MarketStack API")
	}

	return res, nil
}

func (rest *MarketStackRest) Request(symbol string, from models.DateParam, to models.DateParam) ([]byte, error) {
	rest.mu.RLock()
	now := time.Now().Unix()
	if rest.lockedTo != 0 && rest.lockedTo > now {
		offset := time.Second * time.Duration(rest.lockedTo-now)
		logger.Warn.Printf("Reached API rate limit! Waiting %s.", offset.String())
		time.Sleep(offset)
	}

	request, err := http.NewRequest("GET", fmt.Sprintf(
		"https://financialmodelingprep.com/api/v3/historical-price-full/%%5E%s?from=%s&to=%s&apikey=%s",
		symbol,
		FormatDateParam(from),
		FormatDateParam(to),
		rest.AccessKey,
	), nil)

	rest.mu.RUnlock()

	if err != nil {
		return nil, errors.New("failed to initialize new request: " + err.Error())
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(request)
	if err != nil {
		rest.mu.Lock()
		rest.fails++
		if rest.fails == 3 {
			rest.mu.Unlock()
			logger.Error.Panicln("failed to make http request 3 times in a row (check internet connection and/or app credentials)")
		} else {
			offset := time.Millisecond * time.Duration(250*rest.fails)
			rest.mu.Unlock()
			time.Sleep(offset)
			rest.mu.Lock()
			rest.fails = 0
			rest.mu.Unlock()
			return rest.Request(symbol, from, to) // Try again after potential internet connection failure.
		}
	}
	defer res.Body.Close()

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
