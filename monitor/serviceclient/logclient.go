package serviceclient

type ServiceClient struct {
	Domain string
}

func New(domain string) (*ServiceClient, error) {
	return &ServiceClient{
		Domain: domain,
	}, nil
}
