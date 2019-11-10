package servicemanagement

import (
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/api/servicemanagement/v1"
)

func (a *Api) EnableApi(consumerId, serviceName string) error {
	logrus.Infof("Enabling %s API", serviceName)

	operation, err := a.sm.Services.Enable(serviceName, &servicemanagement.EnableServiceRequest{
		ConsumerId: consumerId,
	}).Do()
	if err != nil {
		return err
	}

	for {
		operation, err := a.sm.Operations.Get(operation.Name).Do()
		if err != nil {
			return err
		}

		if operation.Done {
			logrus.Infof("Enabled %s API", serviceName)
			break
		} else {
			logrus.Infof("Enabling still running. Polling again in 2 Seconds")
			time.Sleep(2 * time.Second)
		}
	}

	return nil
}
