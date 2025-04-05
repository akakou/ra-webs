package logclient

type LogClient struct {
	Domain string
}

func New(domain string) (*LogClient, error) {
	return &LogClient{
		Domain: domain,
	}, nil
}
