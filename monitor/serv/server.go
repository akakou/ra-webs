package serv

import (
	"fmt"

	"github.com/akakou/ctstream/monitor/sslmate"
	golangutils "github.com/akakou/go-utils"
	goutils "github.com/akakou/go-utils"
	"github.com/akakou/ra-webs/monitor"
	browsernotifier "github.com/akakou/ra-webs/monitor/notifier/browser"
	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
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
	ct, err := sslmate.DefaultCTClient("")
	if err != nil {
		return nil, err
	}

	notifier, err := browsernotifier.Default()
	if err != nil {
		return nil, err
	}

	monitor, err := monitor.Default(ct, notifier)
	if err != nil {
		return nil, err
	}
	ct.Domain = monitor.TADomain

	adminToken, err := goutils.RandomHex(RANDOM_SIZE)
	if err != nil {
		return nil, errors.Wrap(err, errFailPickupRandom.Error())
	}
	adminToken = golangutils.GetEnv("RA_WEBS_SERVICE_TOKEN", adminToken)
	fmt.Printf("Admin token is: %s\n", adminToken)

	return New(monitor, adminToken)
}

func (server *MonitorServer) Run(address string, e *echo.Echo) error {
	go server.Monitor.Run()

	return e.Start(address)
}

func (server *MonitorServer) Close() {
	server.Monitor.Close()
}
