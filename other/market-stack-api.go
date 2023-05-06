package other

import (
	"errors"
	"fmt"
	"index-bot/logger"
	"io"
	"net/http"
	"strconv"
	"time"
)

type MarketStackRest struct {
	AccessKey   string
	requestLeft uint16
	lockedTo    int64  // Timestamp (in ms) to when it's locked, 0 means there's no lock.
	fails       uint16 // If request failed, try again up to 3 times (delay 250/500/750ms) - after 3rd failed attempt => panic
}

func (rest *MarketStackRest) Request(route string, jsonPayload interface{}) ([]byte, error) {
	now := time.Now().Unix()
	if rest.lockedTo != 0 && rest.lockedTo > now {
		offset := time.Second * time.Duration(rest.lockedTo-now)
		logger.Warn.Printf("Reached API rate limit! Waiting %s.", offset.String())
		time.Sleep(time.Second * time.Duration(rest.lockedTo-now))
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.marketstack.com/v1%s&access_key=%s", route, rest.AccessKey), nil)
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
