package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"index-bot/db"
	"index-bot/logger"

	tempest "github.com/Amatsagu/Tempest"
)

var BlacklistAdd tempest.Command = tempest.Command{
	Name:          "add",
	Description:   "Puts set server on blacklist. Bot won't work on blacklisted servers.",
	AvailableInDM: false,
	Options: []tempest.Option{
		{
			Name:        "guild-id",
			Description: "The id of server you want to shadow ban.",
			Type:        tempest.OPTION_STRING,
			Required:    true,
		},
	},
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

		value, _ := itx.GetOptionValue("guild-id")
		guildID := tempest.StringToSnowflake(value.(string))
		if guildID == 0 {
			itx.SendLinearReply("Provided value is not a valid guild id.", false)
			return
		}

		if _, err := db.Conn.Exec("INSERT INTO guild (id, channel_id, blocked) VALUES (?, 0, true) ON DUPLICATE KEY UPDATE blocked = true;", guildID); err != nil {
			itx.SendLinearReply("Something went wrong... Please try again later and report issue if it happens again.", false)
			logger.Error.Println(err)
			return
		}

		itx.SendLinearReply(fmt.Sprintf("Successfully banned guild with ID = **%s**.", guildID.String()), false)
	},
}
