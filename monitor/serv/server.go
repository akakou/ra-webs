package serv

import (
	"fmt"

	goutils "github.com/akakou/go-utils"
	golangutils "github.com/akakou/golang-utils"
	"github.com/akakou/ra_webs/monitor"
	"github.com/cockroachdb/errors"
)

var errFailPickupRandom = fmt.Errorf("failed to generate random string")

type MonitorServer struct {
	AdminToken string
	Monitor    *monitor.Monitor
}

func New(monitor *monitor.Monitor, adminToken string) (*MonitorServer, error) {
	return &MonitorServer{
		Monitor:    monitor,
		AdminToken: adminToken,
	}, nil
}

func Default() (*MonitorServer, error) {
	monitor, err := monitor.Default()
	if err != nil {
		return nil, err
	}

	adminToken, err := goutils.RandomHex(RANDOM_SIZE)
	if err != nil {
		return nil, errors.Wrap(err, errFailPickupRandom.Error())
	}
	adminToken = golangutils.GetEnv("ADMIN_TOKEN", adminToken)

	fmt.Printf("Admin token generated: %s\n", adminToken)

	return New(monitor, adminToken)
}
