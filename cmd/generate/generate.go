package create

import (
	. "github.com/fionera/TeamDriveManager/cmd"
	generateRclone "github.com/fionera/TeamDriveManager/cmd/generate/rclone"
	"github.com/urfave/cli"
)

func init() {
	RegisterCommand(
		cli.Command{
			Name:  "generate",
			Usage: "All commands that generate something, like rclone configs",
			Subcommands: []cli.Command{
				generateRclone.NewCommand(),
			},
		})
}
