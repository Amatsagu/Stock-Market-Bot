package commands

import (
	"database/sql"
	"errors"
	"index-bot/db"
	"index-bot/events"
	"index-bot/logger"

	tempest "github.com/Amatsagu/Tempest"
)

var SwitchCrons tempest.Command = tempest.Command{
	Name:          "switch-crons",
	Description:   "Enabled/Disables automatic market index changes from posting.",
	AvailableInDM: false,
	SlashCommandHandler: func(itx tempest.CommandInteraction) {
		err := db.Conn.QueryRow("SELECT id FROM operator WHERE id = ?;", itx.Member.User.ID).Scan(&itx.Member.User.ID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				itx.SendLinearReply("This command can only be used by MarketMornings staff members.", false)
				return
			} else {
				itx.SendLinearReply("Something went wrong... Please try again later and report issue if it happens again.", false)
				logger.Error.Println(err)
				return
			}
		}

		enabled := events.AllowCrons.Load()
		if enabled {
			events.AllowCrons.Store(false)
			logger.Info.Printf("Administrator %s (%s) disabled cron jobs!", itx.Member.User.Username, itx.Member.User.ID.String())
			itx.SendLinearReply("Successfully disabled all scheduled cron jobs. Remember to run this command again if you want to again start receiving automatic market change posts.", false)
			return
		} else {
			events.AllowCrons.Store(true)
			logger.Info.Printf("Administrator %s (%s) enabled cron jobs!", itx.Member.User.Username, itx.Member.User.ID.String())
			itx.SendLinearReply("Successfully scheduled all cron jobs.", false)
			return
		}
	},
}
