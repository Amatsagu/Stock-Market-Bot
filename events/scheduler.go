package events

import (
	"sync/atomic"
	"time"

	tempest "github.com/Amatsagu/Tempest"
)

var AllowCrons atomic.Bool = atomic.Bool{}

func InitCronJobs(client *tempest.Client) {
	AllowCrons.Store(true)
	go loopCron(client, "05:00", FirstJob)
	go loopCron(client, "14:00", FirstJob)
	go loopCron(client, "19:00", SecondJob)
	go loopCron(client, "23:00", ThirdJob)
}

func loopCron(client *tempest.Client, ts string, fn func(client *tempest.Client)) {
	startTime := nextDailyCronTime(ts)
	sleepUntil := time.Until(startTime)

	time.Sleep(sleepUntil)
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		enabled := AllowCrons.Load()
		if enabled {
			fn(client)
		}
		<-ticker.C
	}
}

func nextDailyCronTime(ts string) time.Time {
	t, err := time.Parse("15:04", ts)
	if err != nil {
		panic(err)
	}

	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, now.Location())

	if now.After(next) {
		next = next.Add(24 * time.Hour)
	}

	return next
}
