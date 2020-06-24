package generate

import (
	"github.com/urfave/cli"

	. "github.com/fionera/TeamDriveManager/cmd"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "generate",
			Usage: "All commands that generate something, like rclone configs",
			Subcommands: []cli.Command{
				NewGenerateRcloneCommand(),
			},
		})
}
