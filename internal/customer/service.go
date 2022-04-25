package customer

// NewService for initialize service
func NewService(repo Repo) Service {
	return &service{
		repo: repo,
	}
}

// Service will contain all the function that can be used by service
type Service interface {
	PutEditCustomer(body EditCustomerRequest) error
}

type service struct {
	repo Repo
}

func (s service) PutEditCustomer(body EditCustomerRequest) error {
	return nil
}
