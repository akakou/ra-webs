package monitor

func (monitor *Monitor) Run() {
	monitor.CT.Run(monitor)
}
