package serviceaccount

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"

	"github.com/fionera/TeamDriveManager/api"
	"github.com/fionera/TeamDriveManager/api/drive"
	. "github.com/fionera/TeamDriveManager/config"
)

func NewCommand() cli.Command {
	return cli.Command{
		Name:   "rclone",
		Usage:  "Generate a rclone config",
		Action: CmdCreateProject,
		Flags:  []cli.Flag{},
	}
}

func CmdCreateProject(c *cli.Context) {
	filter := strings.Join(c.Args(), " ")

	if filter != "" {
		logrus.Infof("Using filter `%s`", filter)
	}

	client, err := api.CreateClient(App.AppConfig.ServiceAccountFile, App.AppConfig.Impersonate, []string{drive.DriveScope})
	if err != nil {
		logrus.Error(err)
		return
	}

	driveApi, err := drive.NewApi(client)
	if err != nil {
		logrus.Error(err)
		return
	}

	boolResponse := false
	confirm := &survey.Confirm{
		Message: "Use Domain Admin access?",
		Default: false,
	}

	err = survey.AskOne(confirm, &boolResponse, nil)
	if err != nil {
		logrus.Panic(err)
		return
	}

	var list = driveApi.ListTeamDrives
	if boolResponse {
		list = driveApi.ListAllTeamDrives
	}

	teamDrives, err := list()
	if err != nil {
		logrus.Panic(err)
		return
	}

	sb := strings.Builder{}
	for _, teamDrive := range teamDrives {
		if !strings.HasPrefix(teamDrive.Name, filter) {
			continue
		}

		sb.WriteString(fmt.Sprintf("[%s]\n", strings.NewReplacer("/", "_", " ", "").Replace(teamDrive.Name)))
		sb.WriteString("type = drive\n")
		sb.WriteString("scope = drive\n")
		sb.WriteString(fmt.Sprintf("teamdrive_id = %s\n", teamDrive.Id))
		sb.WriteString("\n")
	}

	fmt.Println(sb.String())
	_ = ioutil.WriteFile("rclone.conf", []byte(sb.String()), 0644)
}
