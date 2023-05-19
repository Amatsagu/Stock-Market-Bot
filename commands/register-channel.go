package commands

import (
	"fmt"
	"index-bot/db"
	"index-bot/logger"

	tempest "github.com/Amatsagu/Tempest"
)

var RegisterChannel tempest.Command = tempest.Command{
	Name:          "register-channel",
	Description:   "Binds selected channel to be used to broadcast market changes.",
	AvailableInDM: false,
	Options: []tempest.Option{
		{
			Name:        "channel",
			Description: "Text channel you want to use for broadcasting.",
			Type:        tempest.OPTION_CHANNEL,
			Required:    true,
			ChannelTypes: []tempest.ChannelType{
				tempest.CHANNEL_GUILD_TEXT,
			},
		},
	},
	SlashCommandHandler: func(itx tempest.CommandInteraction) {
		if itx.Member.PermissionFlags&8 != 8 {
			itx.SendLinearReply("This command can only be used by members with `administrator` permission.", false)
			return
		}

		value, _ := itx.GetOptionValue("channel")
		channelID := tempest.StringToSnowflake(value.(string))

		if _, err := db.Conn.Exec("INSERT INTO guild (id, channel_id) VALUES (?, ?) ON DUPLICATE KEY UPDATE channel_id = ?;", itx.GuildID, channelID, channelID); err != nil {
			itx.SendLinearReply("Something went wrong... Please try again later and report issue if it happens again.", false)
			logger.Error.Println(err)
			return
		}

		itx.SendLinearReply(fmt.Sprintf("Successfully set <#%s> as new broadcast channel.", channelID.String()), false)
	},
}
