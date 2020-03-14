package combinations

import (
	"github.com/urfave/cli"

	. "github.com/fionera/TeamDriveManager/cmd"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "combinations",
			Usage: "Commands that do multiple things",
			Subcommands: []cli.Command{
				NewProjectAccountsKeysCommand(),
				NewRegenerateKeysCommand(),
			},
		})
}
