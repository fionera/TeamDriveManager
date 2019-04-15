package group

import (
	"github.com/codegangsta/cli"
)

func NewAssignGroupCmd() cli.Command {
	return cli.Command{
		Name:   "group",
		Usage:  "Assign all users from the config to the corresponding TeamDrives by using Groups",
		Action: CmdAssignGroup,
		Flags:  []cli.Flag{},
	}
}

func CmdAssignGroup() {
	//TODO
}
