package serviceclient

type ServiceClient struct {
	BaseURL string
}

func New(base string) (*ServiceClient, error) {
	return &ServiceClient{
		BaseURL: base,
	}, nil
}

func (sc *ServiceClient) SetDomain(baseURL string) {
	sc.BaseURL = baseURL
}
