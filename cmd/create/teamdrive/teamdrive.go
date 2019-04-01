package teamdrive

import (
	"github.com/codegangsta/cli"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:     "teamdrive",
		Usage:    "Create a Teamdrive",
		Action:   CmdCreateTeamdrive,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name: "service_account_file",
			},
		},
	}
}

func CmdCreateTeamdrive(c *cli.Context) {

}
