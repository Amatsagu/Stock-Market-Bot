package events

import (
	"index-bot/db"
	"index-bot/logger"
	"index-bot/other"
	"strconv"
	"time"

	tempest "github.com/Amatsagu/Tempest"
)

// 7AM CST Year-to-date returns Jan1st-to-current day of data - Nasdaq Composite Index (IXIC.INDX), S&P 500 Index (GSPC.INDX), Dow Jones Industrial Average (DJI.INDX).
// 10PM CST Year-to-date returns   Jan1st-to-current day of data - Nasdaq Composite Index (IXIC.INDX), S&P 500 Index (GSPC.INDX), Dow Jones Industrial Average (DJI.INDX)
func FirstJob(client *tempest.Client) {
	logger.Info.Println("Performing first cron job...")
	channels, err := db.SelectAllValidChannels(nil)
	if err != nil {
		logger.Error.Println("Failed to fetch all valid text channels from database! Terminating first cron job. Error:", err)
		return
	}

	filename := strconv.FormatInt(time.Now().UnixMilli(), 36) + ".png"
	err = other.DrawUSYTDLeaderboard(&db.MarketRest, "./out/"+filename)
	if err != nil {
		logger.Error.Println("Failed to render leaderboard image! Terminating first cron job. Error:", err)
		return
	}

	var success int = 0
	for _, ID := range channels {
		_, err = client.SendLinearMessage(ID, "https://vironicer.kikuri-bot.com/"+filename)
		if err == nil {
			success++
		}
		time.Sleep(time.Millisecond * 250)
	}

	logger.Info.Printf("Finished first cron job! Successfuly sent leaderboard image URL to %d/%d text channels (%d failed).\n", success, len(channels), len(channels)-success)
}

// 12PM CST Overseas year-to-date returns Jan1st-to-current day of data - DAX Performance Index (CDAXX.INDX), FTSE 100 Index (UKX.INDX), CAC 40 (FCHI.INDX), SSE Composite (SSE180.INDX), Hang Seng Index (HSI.INDX), Nikkei 225 (N225.INDX)
func SecondJob(client *tempest.Client) {
	logger.Info.Println("Performing second cron job...")
	channels, err := db.SelectAllValidChannels(nil)
	if err != nil {
		logger.Error.Println("Failed to fetch all valid text channels from database! Terminating second cron job. Error:", err)
		return
	}

	filename := strconv.FormatInt(time.Now().UnixMilli(), 36) + ".png"
	err = other.DrawInternationalYTDLeaderboard(&db.MarketRest, "./out/"+filename)
	if err != nil {
		logger.Error.Println("Failed to render leaderboard image! Terminating second cron job. Error:", err)
		return
	}

	var success int = 0
	for _, ID := range channels {
		_, err = client.SendLinearMessage(ID, "https://vironicer.kikuri-bot.com/"+filename)
		if err == nil {
			success++
		}
		time.Sleep(time.Millisecond * 250)
	}

	logger.Info.Printf("Finished second cron job! Successfuly sent leaderboard image URL to %d/%d text channels (%d failed).\n", success, len(channels), len(channels)-success)
}

// 4PM CST Todays changes End of day data - Nasdaq Composite Index (IXIC.INDX), S&P 500 Index (GSPC.INDX), Dow Jones Industrial Average (DJI.INDX)
func ThirdJob(client *tempest.Client) {
	logger.Info.Println("Performing third cron job...")
	channels, err := db.SelectAllValidChannels(nil)
	if err != nil {
		logger.Error.Println("Failed to fetch all valid text channels from database! Terminating third cron job. Error:", err)
		return
	}

	filename := strconv.FormatInt(time.Now().UnixMilli(), 36) + ".png"
	err = other.DrawUSEODLeaderboard(&db.MarketRest, "./out/"+filename)
	if err != nil {
		logger.Error.Println("Failed to render leaderboard image! Terminating third cron job. Error:", err)
		return
	}

	var success int = 0
	for _, ID := range channels {
		_, err = client.SendLinearMessage(ID, "https://vironicer.kikuri-bot.com/"+filename)
		if err == nil {
			success++
		}
		time.Sleep(time.Millisecond * 250)
	}

	logger.Info.Printf("Finished third cron job! Successfuly sent leaderboard image URL to %d/%d text channels (%d failed).\n", success, len(channels), len(channels)-success)
}
