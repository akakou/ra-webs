package monitor

type CTMonitor interface {
	Setup(monitor *Monitor) error
	Run() error
}
