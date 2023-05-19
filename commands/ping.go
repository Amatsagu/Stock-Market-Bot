package commands

import (
	"fmt"

	tempest "github.com/Amatsagu/Tempest"
)

var Ping tempest.Command = tempest.Command{
	Name:          "ping",
	Description:   "Sends test command.",
	AvailableInDM: false,
	SlashCommandHandler: func(itx tempest.CommandInteraction) {
		itx.SendLinearReply(fmt.Sprintf("Ping! Took **%dms**.", itx.Client.Ping().Milliseconds()), false)
	},
}
