package combinations

import (
	. "github.com/fionera/TeamDriveManager/cmd"
	"github.com/urfave/cli"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "combinations",
			Usage: "Commands that do multiple things",
			Subcommands: []cli.Command{
				NewProjectAccountsKeysCommand(),
			},
		})
}
